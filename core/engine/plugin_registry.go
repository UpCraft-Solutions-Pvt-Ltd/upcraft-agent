package engine

import "sync"

var (
	pluginRegistryMu sync.RWMutex
	pluginRegistry   = map[string]interface{}{}
)

// RegisterPlugin allows plugins to self-register during init().
func RegisterPlugin(name string, implementation interface{}) {
	pluginRegistryMu.Lock()
	defer pluginRegistryMu.Unlock()
	pluginRegistry[name] = implementation
}

// GetPlugin returns one registered plugin implementation by name.
func GetPlugin(name string) (interface{}, bool) {
	pluginRegistryMu.RLock()
	defer pluginRegistryMu.RUnlock()
	v, ok := pluginRegistry[name]
	return v, ok
}

// ListPlugins returns a copy of the current plugin registry.
func ListPlugins() map[string]interface{} {
	pluginRegistryMu.RLock()
	defer pluginRegistryMu.RUnlock()
	out := make(map[string]interface{}, len(pluginRegistry))
	for k, v := range pluginRegistry {
		out[k] = v
	}
	return out
}
