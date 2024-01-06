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
	s := &coreService{}
	s.BasicService = services.NewBasicService(s.start, s.running, s.stop)
	return s, nil
}

func (s *coreService) start(_ context.Context) error {
	serv, err := Initialize()
	if err != nil {
		return err
	}
	s.server = serv
	return s.server.Init()
}

func (s *coreService) running(_ context.Context) error {
	return s.server.Run()
}

func (s *coreService) stop(failureReason error) error {
	return s.server.Shutdown(context.Background(), failureReason.Error())
}
