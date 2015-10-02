package data

import "os"
import "time"
import "log"
import "sync"
import "encoding/json"
import "github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
import "github.com/coreos/etcd/client"

const ETCDDataSourceKey = "/chaostesting/datasources"

type ETCDProvider struct {
	*sync.RWMutex
	*log.Logger

	Interface
	keys client.KeysAPI

	watcher client.Watcher
	done    chan bool
}

func NewETCDClient(endpoints ...string) (client.KeysAPI, error) {
	c, e := client.New(client.Config{
		Endpoints:               endpoints,
		HeaderTimeoutPerRequest: 2 * time.Second,
	})
	if e != nil {
		return nil, e
	}

	return client.NewKeysAPI(c), nil
}

func NewETCDProvider(l *log.Logger, endpoints ...string) (*ETCDProvider, error) {
	if l == nil {
		l = log.New(os.Stdout, "ETCDProvider ", log.LstdFlags)
	}

	etcd, e := NewETCDClient(endpoints...)
	if e != nil {
		return nil, e
	}

	// reads initial configuration
	provider := &ETCDProvider{&sync.RWMutex{}, l, nil, etcd, nil, make(chan bool)}
	if e := provider.reconfigure(); e != nil {
		return nil, e
	}

	// watch for configuration changes
	watcher := provider.keys.Watcher(ETCDDataSourceKey, nil)
	go func() {
		for {
			select {
			case <-provider.done:
				return
			default:
				resp, e := watcher.Next(context.Background())
				if e != nil {
					provider.Println("watcher error:", e)
				}

				provider.Println("etcd event:", resp.Action, resp.Node.Key)
				if resp.Node.Key == ETCDDataSourceKey {
					if e := provider.reconfigure(); e != nil {
						provider.Println("failed to reconfigure:", e)
					}
				}
			} // select
		} // for
	}()

	// TODO: Start watcher.
	return provider, nil
}

func (p *ETCDProvider) Get(result interface{}, query string, args ...interface{}) error {
	return p.try(func() error {
		return p.Interface.Get(result, query, args...)
	})
}

func (p *ETCDProvider) Select(result interface{}, query string, args ...interface{}) error {
	return p.try(func() error {
		return p.Interface.Select(result, query, args...)
	})
}

func (p *ETCDProvider) Exec(query string, args ...interface{}) error {
	return p.try(func() error {
		return p.Interface.Exec(query, args...)
	})
}

func (p *ETCDProvider) try(action func() error) error {
	p.RLock()
	defer p.RUnlock()

	if p.Interface == nil {
		return ErrNoSources
	}

	return action()
}

func (p *ETCDProvider) reconfigure() error {
	p.Println("reading data sources list from ETCD...")

	resp, e := p.keys.Get(context.Background(), ETCDDataSourceKey, nil)
	if e != nil {
		return e
	}

	sources := []*Source{}
	if e := json.Unmarshal([]byte(resp.Node.Value), &sources); e != nil {
		return e
	}

	p.Lock()
	defer p.Unlock()

	if p.Interface != nil {
		if e := p.Interface.Close(); e != nil {
			p.Println("ignored Close() error:", e)
		}

		p.Interface = nil
	}

	p.Interface = NewBalancer(p.Logger, sources...)
	for _, source := range sources {
		p.Println("*", source.String(), source.DSN)
	}

	return nil
}
