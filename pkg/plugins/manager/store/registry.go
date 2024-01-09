package store

import (
	"context"
	"github.com/xquare-dashboard/pkg/plugins"
)

// Service is responsible for the internal storing and retrieval of plugins.
type Service interface {
	// Plugin finds a plugin by its ID.
	Plugin(ctx context.Context, id string) (*plugins.Plugin, bool)
}
