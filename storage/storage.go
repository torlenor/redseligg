package storage

type backend interface {
	StorePluginData(botID, pluginID, identifier string, data interface{}) error
	GetPluginData(botID, pluginID, identifier string) (interface{}, error)
}
