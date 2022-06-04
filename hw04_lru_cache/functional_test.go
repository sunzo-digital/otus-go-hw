package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//  Тесты на взаимодействие кэша и очереди

func TestPurgeOne(t *testing.T) {
	l := NewList()
	c := lruCache{capacity: 3, queue: l, items: make(map[Key]*ListItem, 3)}

	c.Set(Key("1 key"), "1 val")
	c.Set(Key("2 key"), "2 val")
	c.Set(Key("3 key"), "3 val")
	c.Set(Key("4 key"), "4 val")

	require.True(t, l.Len() == 3)
	require.True(t, l.Front().Value.(cacheItem).value == "4 val")
	require.True(t, l.Back().Value.(cacheItem).value == "2 val")
	require.True(t, l.Front().Next.Value.(cacheItem).value == "3 val" && l.Back().Prev.Value.(cacheItem).value == "3 val")
	require.True(t, l.Front().Prev == nil && l.Back().Next == nil)
}

func TestPurgeAfterMixing(t *testing.T) {
	l := NewList()
	c := lruCache{capacity: 3, queue: l, items: make(map[Key]*ListItem, 3)}

	c.Set(Key("1 key"), "1 val")
	c.Set(Key("2 key"), "2 val")
	c.Set(Key("3 key"), "3 val") // 3 2 1
	c.Set(Key("1 key"), "1 val") // 1 3 2
	c.Get("2 key")               // 2 1 3
	c.Set(Key("4 key"), "4 val") // 4 2 1

	require.True(t, l.Len() == 3)
	require.True(t, l.Front().Value.(cacheItem).value == "4 val")
	require.True(t, l.Back().Value.(cacheItem).value == "1 val")
	require.True(t, l.Front().Next.Value.(cacheItem).value == "2 val" && l.Back().Prev.Value.(cacheItem).value == "2 val")
	require.True(t, l.Front().Prev == nil && l.Back().Next == nil)
}
