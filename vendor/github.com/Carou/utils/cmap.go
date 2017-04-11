package utils

import (
	"encoding/json"
	"hash/fnv"
	"reflect"
	"sync"
)

var (
	shardCount = 32
)

// ConcurrentMap represents concurrent map we are implementing
type ConcurrentMap []*ConcurrentMapShared

// ConcurrentMapShared is a thread safe map for string -> interface{}
type ConcurrentMapShared struct {
	sync.RWMutex
	items map[string]interface{}
}

// NewConcurrentMap returns a new instance of concurrent map
func NewConcurrentMap() ConcurrentMap {
	m := make(ConcurrentMap, shardCount)
	for i := 0; i < shardCount; i++ {
		m[i] = &ConcurrentMapShared{
			items: make(map[string]interface{}),
		}
	}
	return m
}

func (cm ConcurrentMap) getShard(key string) *ConcurrentMapShared {
	hasher := fnv.New32()
	_, err := hasher.Write([]byte(key))
	if err != nil {
		return nil
	}
	return cm[uint(hasher.Sum32())%uint(shardCount)]
}

// Set set a new k-v pair in concurrent map
func (cm *ConcurrentMap) Set(key string, value interface{}) {
	shard := cm.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.items[key] = value
}

// Get gets value based on given key in concurrent map
func (cm ConcurrentMap) Get(key string) (interface{}, bool) {
	shard := cm.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	val, ok := shard.items[key]
	return val, ok
}

// GetElseSet value based on the given key, or set and return the new value if not exist
// If the key already exists, return the existing value and true
// Otherwise return the newValue and false
func (cm ConcurrentMap) GetElseSet(key string, newValue interface{}) (interface{}, bool) {
	shard := cm.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	if val, found := shard.items[key]; found {
		return val, true
	}

	shard.items[key] = newValue
	return newValue, false
}

// Count returns the count of all items in concurrent map
func (cm ConcurrentMap) Count() int {
	count := 0
	for i := 0; i < shardCount; i++ {
		shard := cm[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Has returns true if given key is in concurrent map, else false
func (cm ConcurrentMap) Has(key string) bool {
	shard := cm.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	_, ok := shard.items[key]
	return ok
}

// Remove removes give key in concurrent map
func (cm ConcurrentMap) Remove(key string) {
	shard := cm.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	delete(shard.items, key)
}

// RemoveIfValue removes key ONLY when value is found
func (cm ConcurrentMap) RemoveIfValue(key string, value interface{}) (ret bool, err error) {
	defer func() {
		if e := recover(); e != nil {
			err, _ = e.(error)
		}
	}()
	shard := cm.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	if val, ok := shard.items[key]; ok {
		if reflect.DeepEqual(val, value) {
			delete(shard.items, key)
			ret = true
		}
	}
	return
}

// IsEmpty returns true if concurrent map is empty, else false
func (cm ConcurrentMap) IsEmpty() bool {
	return cm.Count() == 0
}

// Tuple used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type Tuple struct {
	Key string
	Val interface{}
}

// Iter returns an iterator which could be used in a for range loop.
func (cm ConcurrentMap) Iter() []Tuple {
	result := make([]Tuple, 0, cm.Count())
	// Foreach shard.
	for _, shard := range cm {
		// Foreach key, value pair.
		shard.RLock()
		for key, val := range shard.items {
			result = append(result, Tuple{key, val})
		}
		shard.RUnlock()
	}
	return result
}

// MarshalJSON reviles ConcurrentMap "private" variables to json marshal.
func (cm ConcurrentMap) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[string]interface{})

	// Insert items to temporary map.
	for _, item := range cm.Iter() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}
