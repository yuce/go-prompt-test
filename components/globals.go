package components

import "sync"

var globalKeyValues = map[string]interface{}{}
var globalKeyValuesMu = &sync.RWMutex{}

func FindGlobalValue(key string) (interface{}, bool) {
	globalKeyValuesMu.RLock()
	defer globalKeyValuesMu.RUnlock()
	value, ok := globalKeyValues[key]
	return value, ok
}

func UpdateGlobal(key string, value interface{}) {
	globalKeyValuesMu.Lock()
	defer globalKeyValuesMu.Unlock()
	globalKeyValues[key] = value
}
