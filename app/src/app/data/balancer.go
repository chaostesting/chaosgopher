package data

import "os"
import "log"
import "sync"
import "github.com/jmoiron/sqlx"

const MaxRetry = 3

// Balancer is a provider that automatically creates and maintains connections to multiple
// data sources and automatically load balances calls between them.
type Balancer struct {
	*sync.RWMutex
	*log.Logger

	sources     []*Source
	connections map[*Source]*sqlx.DB
}

func NewBalancer(l *log.Logger, sources ...*Source) *Balancer {
	if l == nil {
		l = log.New(os.Stdout, "Balancer ", log.LstdFlags)
	}

	return &Balancer{&sync.RWMutex{}, l, sources, map[*Source]*sqlx.DB{}}
}

func (balancer *Balancer) Get(result interface{}, query string, args ...interface{}) error {
	q, e := balancer.db()
	if e != nil {
		return e
	}

	return sqlx.Get(q, result, query, args...)
}

func (balancer *Balancer) Select(result interface{}, query string, args ...interface{}) error {
	q, e := balancer.db()
	if e != nil {
		return e
	}

	return sqlx.Select(q, result, query, args...)
}

func (balancer *Balancer) Exec(query string, args ...interface{}) error {
	x, e := balancer.db()
	if e != nil {
		return e
	}

	_, e = x.Exec(query, args...)
	return e
}

func (balancer *Balancer) Close() error {
	balancer.Lock()
	defer balancer.Unlock()

	for dsn, conn := range balancer.connections {
		go func() {
			if e := conn.Close(); e != nil {
				balancer.Println("ignored Close() failure:", dsn)
			}
		}()
	}

	balancer.connections = map[*Source]*sqlx.DB{}
	return nil
}

func (balancer *Balancer) db() (*sqlx.DB, error) {
	return balancer._db(0, nil)
}

func (balancer *Balancer) _db(tries int, lastError error) (*sqlx.DB, error) {
	if tries >= MaxRetry {
		if lastError == nil {
			lastError = ErrNoGoodSource
		}

		return nil, lastError
	}

	source := SelectSource(balancer.sources)
	if source == nil {
		return nil, ErrNoSources
	}

	balancer.RLock()
	conn, ok := balancer.connections[source]
	balancer.RUnlock()
	if ok {
		if e := conn.Ping(); e != nil { // connection no longer valid.
			balancer.Println("bad connection, pruning:", source.Name)

			balancer.Lock()
			delete(balancer.connections, source)
			balancer.Unlock()

			return balancer._db(tries, e) // stale connection does not count as retry
		}

		return conn, nil
	}

	balancer.Lock()
	conn, e := sqlx.Connect(source.Driver, source.DSN)
	if e != nil {
		balancer.Unlock()
		return balancer._db(tries+1, e)
	}

	balancer.connections[source] = conn
	balancer.Unlock()
	return conn, nil
}
