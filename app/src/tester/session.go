package main

import "log"
import "time"
import "strings"
import "encoding/json"
import "math/rand"
import "os"
import "app/config"
import "app/data"
import "app/client"
import "github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
import etcdclient "github.com/coreos/etcd/client"

func init() {
	rand.Seed(time.Now().Unix())
}

type Session struct {
	*log.Logger
	controllers map[string]*Controller
	sources     []*data.Source

	etcd  etcdclient.KeysAPI
	api   *client.APIClient
	todos []*data.TodoItem
}

func NewSession(l *log.Logger, containers ...string) (*Session, error) {
	if l == nil {
		l = log.New(os.Stderr, "[session] ", log.LstdFlags)
	}

	etcd, e := data.NewETCDClient(config.Get().ETCDEndpoint)
	if e != nil {
		return nil, e
	}

	// initial configuration becomes "known defaults"
	resp, e := etcd.Get(context.Background(), data.ETCDDataSourceKey, nil)
	if e != nil {
		return nil, e
	}

	sources := []*data.Source{}
	if e := json.Unmarshal([]byte(resp.Node.Value), &sources); e != nil {
		return nil, e
	}

	client := &client.APIClient{Host: config.Get().APIEndpoint}
	controllers := map[string]*Controller{}
	for _, container := range containers {
		controller, e := NewController(container, "")
		if e != nil {
			return nil, e
		}

		controllers[container] = controller
	}

	// special filler controller
	controller, e := NewController("filler", "hyperworks/disk-filler")
	if e != nil {
		return nil, e
	}

	controllers["filler"] = controller

	return &Session{l, controllers, sources, etcd, client, nil}, nil
}

func (s *Session) Query() error {
	s.Println("querying TODOs:")
	todos, e := s.api.GetAllTodos()
	if e != nil {
		return e
	}

	for _, todo := range todos {
		desc := todo.Description
		if todo.Completed {
			desc = strings.Repeat("-", len(desc))
		}

		s.Printf(" %d. %s", todo.ID, desc)
	}

	s.todos = todos
	return nil
}

func (s *Session) Mutate() error {
	s.Println("toggle random TODOs")

	for _, todo := range s.todos {
		if rand.Int()%2 == 0 {
			s.Println("*", todo.Description)
			todo.Completed = !todo.Completed

			_, e := s.api.PatchTodo(todo.ID, todo)
			if e != nil {
				return e
			}
		}
	}

	return nil
}

func (s *Session) Create(container string) error {
	s.Println("creating", container)
	return s.controllers[container].Create()
}

func (s *Session) Destroy(container string) error {
	s.Println("destroying", container)
	return s.controllers[container].Destroy()
}

func (s *Session) Start(container string) error {
	s.Println("starting", container)
	return s.controllers[container].Start()
}

func (s *Session) Stop(container string) error {
	s.Println("stopping", container)
	return s.controllers[container].Stop()
}

func (s *Session) Status(containers ...string) error {
	s.Println("checking container status")
	for _, name := range containers {
		controller := s.controllers[name]
		state, e := controller.Started()
		if e != nil {
			return e
		}

		if state {
			s.Println("*", name, "is running")
		} else {
			s.Println("*", name, "stopped")
		}
	}

	return nil
}

func (s *Session) EnsureStarted(containers ...string) error {
	s.Println("ensure container started:", containers)
	for _, name := range containers {
		controller := s.controllers[name]
		state, e := controller.Started()
		if e != nil {
			return e
		}

		if !state {
			if e := controller.Start(); e != nil {
				return e
			}

			s.Println("*", name, "started.")
		}
	}

	return nil
}

func (s *Session) Reconfigure(instances ...string) error {
	valid := map[string]bool{}
	for _, instance := range instances {
		valid[instance] = true
	}

	s.Println("reconfiguring application...")
	sources := []*data.Source{}
	for _, source := range s.sources {
		if v, ok := valid[source.Name]; ok && v {
			s.Println("*", source)
			sources = append(sources, source)
		}
	}

	bytes, e := json.Marshal(sources)
	if e != nil {
		return e
	}

	_, e = s.etcd.Set(context.Background(), data.ETCDDataSourceKey, string(bytes), nil)
	if e != nil {
		return e
	}

	return nil
}
