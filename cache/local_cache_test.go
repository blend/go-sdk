package cache

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
)

var (
	_ graceful.Graceful = (*LocalCache)(nil)
)

type itemKey struct{}
type altItemKey struct{}

func TestLocalCache(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache()

	t1 := time.Date(2019, 06, 14, 12, 10, 9, 8, time.UTC)
	t2 := time.Date(2019, 06, 14, 00, 01, 02, 03, time.UTC)
	t3 := time.Date(2019, 06, 14, 12, 01, 02, 03, time.UTC)

	c.Set(itemKey{}, "foo", OptValueExpires(t1))
	assert.True(c.Has(itemKey{}))
	assert.False(c.Has(altItemKey{}))

	found, ok := c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("foo", found)

	c.Set(altItemKey{}, "alt-bar")
	assert.True(c.Has(itemKey{}))
	assert.True(c.Has(altItemKey{}))

	found, ok = c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("foo", found)

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("alt-bar", found)

	c.Set(itemKey{}, "foo-2", OptValueExpires(t2))
	assert.Equal(t2, c.Data[itemKey{}].Expires)
	c.Set(altItemKey{}, "alt-bar-2", OptValueExpires(t3))
	assert.Equal(t3, c.Data[altItemKey{}].Expires)

	found, ok = c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("foo-2", found)

	c.Remove(itemKey{})
	assert.False(c.Has(itemKey{}))
	assert.True(c.Has(altItemKey{}))

	found, ok = c.Get(itemKey{})
	assert.False(ok)
	assert.Nil(found)

	c.Set(itemKey{}, "bar", OptValueExpires(time.Now().UTC().Add(-time.Hour)))
	assert.True(c.Has(itemKey{}))
	assert.True(c.Has(altItemKey{}))
}

func try(action func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	action()
	return
}

func TestLocalCacheKeyPanic(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache()

	assert.NotNil(try(func() {
		c.Set(nil, "bar")
	}))
	assert.NotNil(try(func() {
		c.Set([]int{}, "bar")
	}))
}

func TestLocalCacheGetOrSet(t *testing.T) {
	assert := assert.New(t)

	valueProvider := func() (interface{}, error) { return "foo", nil }

	lc := NewLocalCache()
	found, ok, err := lc.GetOrSet(itemKey{}, valueProvider)
	assert.Nil(err)
	assert.False(ok)
	assert.Equal("foo", found)
	assert.True(lc.Has(itemKey{}))
	assert.Equal(itemKey{}, lc.LRU.Peek().Key)

	found, ok, err = lc.GetOrSet(itemKey{}, valueProvider)
	assert.Nil(err)
	assert.True(ok)
	assert.Equal("foo", found)

	lc.Set(itemKey{}, "bar")

	found, ok, err = lc.GetOrSet(itemKey{}, valueProvider)
	assert.Nil(err)
	assert.True(ok)
	assert.Equal("bar", found)

	lc.Remove(itemKey{})
	assert.False(lc.Has(itemKey{}))

	found, ok, err = lc.GetOrSet(itemKey{}, valueProvider)
	assert.Nil(err)
	assert.False(ok)
	assert.Equal("foo", found)
}

func TestLocalCacheGetOrSetError(t *testing.T) {
	assert := assert.New(t)

	valueProvider := func() (interface{}, error) {
		return nil, fmt.Errorf("test")
	}

	lc := NewLocalCache()

	found, ok, err := lc.GetOrSet("test", valueProvider)
	assert.NotNil(err)
	assert.False(ok)
	assert.Nil(found)

	assert.False(lc.Has("test"))
}

func TestLocalCacheGetOrSetDoubleCheckRace(t *testing.T) {
	assert := assert.New(t)

	didSet := make(chan struct{})
	valueProvider := func() (interface{}, error) {
		<-didSet
		return "foo", nil
	}

	lc := NewLocalCache()

	go func() {
		lc.Set("test", "bar2")
		close(didSet)
	}()

	found, ok, err := lc.GetOrSet("test", valueProvider)
	assert.Nil(err)
	assert.True(ok)
	assert.Equal("bar2", found)
}

func TestLocalCacheSetUpdatesLRU(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache()
	c.Set("k1", "v1", OptValueTTL(0))
	time.Sleep(time.Millisecond)
	c.Set("k2", "v2", OptValueTTL(0))
	assert.Equal("k1", c.LRU.Peek().Key)

	// Should trigger sorting of underlying LRU so k2 can be
	// deleted in next sweep
	time.Sleep(time.Millisecond)
	c.Set("k1", "v3", OptValueTTL(time.Second))

	assert.Equal("k2", c.LRU.Peek().Key)

	c.Sweep(context.Background())
	assert.True(c.Has("k1"))
	assert.False(c.Has("k2"))
}

func TestLocalCacheSweep(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache()

	var didSweep, didRemove bool
	c.Set(itemKey{}, "foo",
		OptValueTimestamp(time.Now().UTC().Add(-2*time.Minute)),
		OptValueTTL(time.Minute),
		OptValueOnRemove(func(_ interface{}, reason RemovalReason) {
			if reason == Expired {
				didSweep = true
			}
		}),
	)
	found, ok := c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("foo", found)

	c.Set(altItemKey{}, "bar",
		OptValueTTL(time.Minute),
	)

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("bar", found)

	c.Sweep(context.Background())

	found, ok = c.Get(itemKey{})
	assert.False(ok)
	assert.Nil(found)
	assert.True(didSweep)
	assert.False(didRemove)

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("bar", found)
}

func TestLocalCacheStartSweeping(t *testing.T) {
	assert := assert.New(t)

	c := NewLocalCache(OptLocalCacheSweepInterval(time.Millisecond))

	didSweep := make(chan struct{})
	c.Set(itemKey{}, "a value",
		OptValueTTL(time.Microsecond),
		OptValueOnRemove(func(_ interface{}, reason RemovalReason) {
			if reason == Expired {
				close(didSweep)
			}
		}),
	)

	found, ok := c.Get(itemKey{})
	assert.True(ok)
	assert.Equal("a value", found)

	c.Set(altItemKey{}, "bar",
		OptValueTTL(time.Minute),
	)

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("bar", found)

	go c.Start()
	<-c.NotifyStarted()
	defer c.Stop()
	<-didSweep

	found, ok = c.Get(itemKey{})
	assert.False(ok)
	assert.Nil(found)

	found, ok = c.Get(altItemKey{})
	assert.True(ok)
	assert.Equal("bar", found)
}

func TestLocalCacheStats(t *testing.T) {
	assert := assert.New(t)

	t1 := time.Date(2019, 06, 14, 12, 10, 9, 8, time.UTC)
	t2 := time.Date(2019, 06, 14, 00, 01, 02, 03, time.UTC)
	t3 := time.Date(2019, 06, 14, 12, 01, 02, 03, time.UTC)

	lc := NewLocalCache()

	lc.Set("foo", "bar", OptValueTimestamp(t1))
	lc.Set("foo2", "bar2", OptValueTimestamp(t2))
	lc.Set("foo3", "bar3", OptValueTimestamp(t3))

	stats := lc.Stats()
	assert.Equal(3, stats.Count)
	assert.Equal(24, stats.SizeBytes)
	assert.NotZero(stats.MaxAge)
}

func BenchmarkLocalCache(b *testing.B) {
	for x := 0; x < b.N; x++ {
		benchLocalCache(1024)
	}
}

func benchLocalCache(items int) {
	lc := NewLocalCache()
	for x := 0; x < items; x++ {
		lc.Set(x, strconv.Itoa(x), OptValueTTL(time.Millisecond))
	}
	for x := 0; x < items; x++ {
		lc.Set(x, strconv.Itoa(x), OptValueTTL(time.Second))
	}
	var value interface{}
	var ok bool
	for x := 0; x < items; x++ {
		value, ok = lc.Get(x)
		if !ok {
			panic("value not found")
		}
		if value.(string) != strconv.Itoa(x) {
			panic("wrong value")
		}
	}
	lc.Sweep(context.Background())
}
