package memorystorage

type memoryStorage map[string]interface{}
type pluginStorage map[string]memoryStorage
type botStorage map[string]pluginStorage
