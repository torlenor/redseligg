package memorybackend

type storage map[string]interface{}
type pluginStorage map[string]storage
type botStorage map[string]pluginStorage
