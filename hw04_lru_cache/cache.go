package main

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	lItem, ok := c.items[key]
	cItem := cacheItem{value: value, key: key}

	if ok {
		lItem.Value = cItem
		c.queue.MoveToFront(lItem)
		return ok
	}

	c.items[key] = c.queue.PushFront(cItem)

	if c.queueIsFull() {
		last := c.queue.Back()
		c.queue.Remove(last)
		delete(c.items, last.Value.(cacheItem).key)
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	lItem, ok := c.items[key]

	if ok {
		c.queue.MoveToFront(lItem)
		return lItem.Value.(cacheItem).value, ok
	}

	return nil, ok
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}

func (c *lruCache) queueIsFull() bool {
	return c.queue.Len() > c.capacity
}
