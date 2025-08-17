package registries

import "sync"

type TypedSyncMap[K comparable, V any] struct {
	m sync.Map
}

func (tm *TypedSyncMap[K, V]) All() []V {
	values := []V{}
	tm.m.Range(func(k any, v any) bool {
		values = append(values, v.(V))
		return true
	})
	return values
}
func (tm *TypedSyncMap[K, V]) Keys() []K {
	keys := []K{}
	tm.m.Range(func(k any, v any) bool {
		keys = append(keys, k.(K))
		return true
	})
	return keys
}

func (tm *TypedSyncMap[K, V]) Store(key K, value V) {
	tm.m.Store(key, value)
}

func (tm *TypedSyncMap[K, V]) Load(key K) (V, bool) {
	val, ok := tm.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return val.(V), true
}

func (tm *TypedSyncMap[K, V]) Delete(key K) {
	tm.m.Delete(key)
}
