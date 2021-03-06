package cache

import (
	"sync"

	"container/list"
)

// lruCacheEntry is a item of lru cache
type lruCacheEntry struct {
	key string
	val interface{}
}

// init init a cache entry
func (ce *lruCacheEntry) init(key string, val interface{}) {
	ce.key = key
	ce.val = val
}

// lruCache is a cacher use lru eliminate algorithm
type lruCache struct {
	cacheData  *list.List
	cacheIndex map[string]*list.Element
	maxSize    int
	lock       *sync.RWMutex
}

// Init init lru cacher
func (lc *lruCache) Init(config string) (err error) {
	lc.cacheData = list.New()
	if lc.maxSize, err = parseMaxSize(config); err == nil {
		lc.cacheIndex = make(map[string]*list.Element, lc.maxSize)
		lc.lock = new(sync.RWMutex)
	}
	return
}

// Init init lru cacher
func (lc *lruCache) InitVals(config string, values map[string]interface{}) (err error) {
	lc.cacheData = list.New()
	if lc.maxSize, err = parseMaxSize(config); err == nil {
		fixSize(values, lc.maxSize)
		lc.cacheIndex = make(map[string]*list.Element, lc.maxSize)
		lc.lock = new(sync.RWMutex)
		for k, v := range values {
			lc.cacheIndex[k] = lc.cacheData.PushFront(v)
		}
	}
	return
}

// Size return current cache count
// it's safe for concurrent
func (lc *lruCache) Size() int {
	lc.lock.RLock()
	size := lc.size()
	lc.lock.RUnlock()
	return size
}

// size is same as Size, but don't require read lock
func (lc *lruCache) size() int {
	return len(lc.cacheIndex)
}

// Cap return cache capacity
func (lc *lruCache) Cap() int {
	// lc.lock.RLock() current it's not need
	c := lc.cap()
	// lc.lock.RUnlock()
	return c
}

// cap is same as Cap, but don't require read lock
func (lc *lruCache) cap() int {
	return lc.maxSize
}

// Get return value of the key, if not exist, nil returned
func (lc *lruCache) Get(key string) (val interface{}) {
	lc.lock.RLock()
	elem, has := lc.cacheIndex[key]
	if has {
		val = elem.Value.(*lruCacheEntry).val
		lc.cacheData.MoveToFront(elem)
	}
	lc.lock.RUnlock()
	return
}

func (lc *lruCache) IsExist(key string) bool {
	lc.lock.RLock()
	_, has := lc.cacheIndex[key]
	lc.lock.RUnlock()
	return has
}

// Remove remove key and it's value from cache
func (lc *lruCache) Remove(key string) {
	lc.lock.Lock()
	elem, has := lc.cacheIndex[key]
	if has {
		lc.cacheData.Remove(elem)
		delete(lc.cacheIndex, key)
	}
	lc.lock.Unlock()
}

// Set add an key-value to cache, if key already exist in cache, update it's value
func (lc *lruCache) Set(key string, val interface{}) {
	lc.set(key, val, true)
}

// Update only update existed key-value, returned value show whether it's successed
func (lc *lruCache) Update(key string, val interface{}) bool {
	return lc.set(key, val, false)
}

// set do actually update cache, the parameter forceSet make a difference when
// key already exist in cache, if forceSet, update it's value, else do nothing
// return value show if operation is successed or not
func (lc *lruCache) set(key string, val interface{}, forceSet bool) (ret bool) {
	var entry *lruCacheEntry
	ret = true
	lc.lock.Lock()
	if elem, has := lc.cacheIndex[key]; !has {
		if !forceSet {
			ret = false
		} else if lc.cap() == lc.size() {
			elem = lc.cacheData.Back() // remove last and reuse entry for new value
			entry = elem.Value.(*lruCacheEntry)
			lc.cacheData.Remove(elem)
			delete(lc.cacheIndex, entry.key)
		} else {
			entry = new(lruCacheEntry)
		}
		entry.init(key, val) // setup value
		lc.cacheIndex[key] = lc.cacheData.PushFront(entry)
	} else {
		elem.Value.(*lruCacheEntry).val = val
		lc.cacheData.MoveToFront(elem)
	}
	lc.lock.Unlock()
	return
}
