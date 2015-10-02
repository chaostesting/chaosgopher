package main

import "log"
import "os"
import "app/config"
import "app/data"

type Services struct {
	Config *config.Config
	Logger *log.Logger
	DB     data.Interface

	Handlers *Handlers
}

func DefaultServices() (*Services, error) {
	cfg := config.Get()
	logger := log.New(os.Stderr, "[chaostesting] ", log.LstdFlags)

	provider, e := data.NewETCDProvider(logger, cfg.ETCDEndpoint)
	if e != nil {
		return nil, e
	}

	srv := &Services{
		Logger: logger,
		DB:     provider,
	}

	srv.Handlers = &Handlers{srv}
	return srv, nil
}
