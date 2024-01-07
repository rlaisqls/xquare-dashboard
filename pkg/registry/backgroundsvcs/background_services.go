package backgroundsvcs

import (
	"github.com/xquare-dashboard/pkg/api"
	"github.com/xquare-dashboard/pkg/registry"
)

func ProvideBackgroundServiceRegistry(
	httpServer *api.HTTPServer,
) *BackgroundServiceRegistry {
	return NewBackgroundServiceRegistry(
		httpServer,
	)
}

// BackgroundServiceRegistry provides background services.
type BackgroundServiceRegistry struct {
	Services []registry.BackgroundService
}

func NewBackgroundServiceRegistry(services ...registry.BackgroundService) *BackgroundServiceRegistry {
	return &BackgroundServiceRegistry{services}
}

func (r *BackgroundServiceRegistry) GetServices() []registry.BackgroundService {
	return r.Services
}
