package server

import (
	"context"
	"github.com/grafana/dskit/services"
)

type coreService struct {
	*services.BasicService
	server *Server
}

func NewService() (*coreService, error) {
	return &coreService{}, nil
}

func (s *coreService) Start() error {
	serv, err := Initialize()
	if err != nil {
		return err
	}
	s.server = serv
	return s.server.Init()
}

func (s *coreService) Running() error {
	return s.server.Run()
}

func (s *coreService) Stop(failureReason error) error {
	return s.server.Shutdown(context.Background(), failureReason.Error())
}
