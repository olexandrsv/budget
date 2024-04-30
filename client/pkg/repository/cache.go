package repository

import (
	"budget/pkg/models"
)

type Cache struct {
	data map[models.Key]any
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[models.Key]any),
	}
}

func cacheGetStore[T models.DataObject[T]](cache *Cache) models.Store[T] {
	var item T
	m, ok := cache.data[item.CacheKey()]
	if !ok {
		return nil
	}

	return m.(models.Store[T])
}

func cachePutStore[T models.DataObject[T]](cache *Cache, store models.Store[T]) {
	var item T
	cache.data[item.CacheKey()] = store
}

func cachePutItem[T models.DataObject[T]](cache *Cache, item T) {
	store, ok := cache.data[item.CacheKey()].(models.Store[T])
	if !ok {
		store = models.NewStore[T]()
		cache.data[item.CacheKey()] = store
	}
	store.Put(item.Index(), item)
}

func cacheGetItem[T models.DataObject[T]](cache *Cache, index int) T {
	var item T
	return cache.data[item.CacheKey()].(models.Store[T]).Get(index)
}

func cacheDeleteItem[T models.DataObject[T]](cache *Cache, index int) {
	var item T
	cache.data[item.CacheKey()].(models.Store[T]).Delete(index)
}

type Index interface {
	Index() int
}
