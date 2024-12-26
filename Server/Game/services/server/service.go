package server

import (
	"Common/Framework/config"
	"context"
	"github.com/hsgames/gold/app"
	"github.com/hsgames/gold/net"
	"github.com/hsgames/gold/net/ws"
	"time"
)

type service struct {
	*ws.Server
}

func (s *service) Name() string {
	return "ws-server"
}

func (s *service) Init() error {
	return nil
}

func (s *service) Start() (err error) {
	return s.ListenAndServe()
}

func (s *service) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	s.Shutdown(ctx)

	return nil
}

func New() (app.Service, error) {
	var (
		s   = &service{}
		err error
	)

	s.Server, err = ws.NewServer(
		s.Name(),
		config.Instance().App.Addr,
		func() net.Handler {
			return &handler{service: s}
		},
		ws.WithBinary(),
	)

	return s, err
}
