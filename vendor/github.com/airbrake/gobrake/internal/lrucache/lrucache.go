package lrucache

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	key     string
	value   interface{}
	addedAt time.Time
}

type Cache struct {
	mu sync.Mutex

	list  *list.List
	table map[string]*list.Element

	maxLen int
}

func New(maxLen int) *Cache {
	return &Cache{
		list:  list.New(),
		table: make(map[string]*list.Element, maxLen),

		maxLen: maxLen,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.get(key)
}

func (c *Cache) get(key string) (interface{}, bool) {
	c.mu.Lock()

	el := c.table[key]
	if el == nil {
		c.mu.Unlock()
		return nil, false
	}

	entry := el.Value.(*entry)
	c.list.MoveToFront(el)
	value := entry.value
	c.mu.Unlock()
	return value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	if el := c.table[key]; el != nil {
		entry := el.Value.(*entry)
		entry.value = value
		c.promote(el, entry)
	} else {
		c.addNew(key, value)
	}
	c.mu.Unlock()
}

func (c *Cache) addNew(key string, value interface{}) {
	newEntry := &entry{
		key:     key,
		value:   value,
		addedAt: time.Now(),
	}
	element := c.list.PushFront(newEntry)
	c.table[key] = element
	c.check()
}

func (c *Cache) promote(el *list.Element, entry *entry) {
	entry.addedAt = time.Now()
	c.list.MoveToFront(el)
}

func (c *Cache) deleteElement(el *list.Element) {
	c.list.Remove(el)
	delete(c.table, el.Value.(*entry).key)
}

func (c *Cache) check() {
	for c.list.Len() > c.maxLen {
		el := c.list.Back()
		c.deleteElement(el)
	}
}
