// Package backendplugin contains backend plugin related logic.
package backendplugin

// PluginFactoryFunc is a function type for creating a BackendPlugin.
type PluginFactoryFunc func(pluginID string) (BackendPlugin, error)
