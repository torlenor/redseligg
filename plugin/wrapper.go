package plugin

import "errors"

// Store wraps the Store call of the Storage API for more convenient use
func (p *AbyleBotterPlugin) Store(identifier string, data interface{}) error {
	if p.Storage == nil {
		return errors.New("No storage API set")
	}
	return p.Storage.Store(p.PluginID, identifier, data)
}

// Get wraps the Get call of the Storage API for more convenient use
func (p *AbyleBotterPlugin) Get(identifier string) (interface{}, error) {
	if p.Storage == nil {
		return nil, errors.New("No storage API set")
	}
	return p.Storage.Get(p.PluginID, identifier)
}
