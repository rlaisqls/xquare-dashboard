package apiregistry

import (
	"context"

	"github.com/xquare-dashboard/pkg/registry"
	"github.com/xquare-dashboard/pkg/registry/apis/example"
	"github.com/xquare-dashboard/pkg/registry/apis/folders"
	"github.com/xquare-dashboard/pkg/registry/apis/playlist"
)

var (
	_ registry.BackgroundService = (*Service)(nil)
)

type Service struct{}

// ProvideRegistryServiceSink is an entry point for each service that will force initialization
// and give each builder the chance to register itself with the main server
func ProvideRegistryServiceSink(
	_ *playlist.PlaylistAPIBuilder,
	_ *example.TestingAPIBuilder,
	_ *folders.FolderAPIBuilder,
) *Service {
	return &Service{}
}

func (s *Service) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}
