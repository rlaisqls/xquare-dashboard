package usagestatssvcs

import (
	"github.com/xquare-dashboard/pkg/registry"
	"github.com/xquare-dashboard/pkg/services/accesscontrol"
	"github.com/xquare-dashboard/pkg/services/user"
)

func ProvideUsageStatsProvidersRegistry(
	accesscontrol accesscontrol.Service,
	user user.Service,
) *UsageStatsProvidersRegistry {
	return NewUsageStatsProvidersRegistry(
		accesscontrol,
		user,
	)
}

type UsageStatsProvidersRegistry struct {
	Services []registry.ProvidesUsageStats
}

func NewUsageStatsProvidersRegistry(services ...registry.ProvidesUsageStats) *UsageStatsProvidersRegistry {
	return &UsageStatsProvidersRegistry{services}
}

func (r *UsageStatsProvidersRegistry) GetServices() []registry.ProvidesUsageStats {
	return r.Services
}
