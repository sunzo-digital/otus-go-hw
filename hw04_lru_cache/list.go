package main

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v, Next: nil, Prev: nil}

	if l.len == 0 {
		l.front = newItem
		l.back = newItem
		l.len++

		return newItem
	}

	next := l.Front()
	newItem.Next = next
	next.Prev = newItem
	l.front = newItem
	l.len++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v, Next: nil, Prev: nil}

	if l.len == 0 {
		l.front = newItem
		l.back = newItem
		l.len++

		return newItem
	}

	prev := l.Back()
	newItem.Prev = prev
	prev.Next = newItem

	l.back = newItem
	l.len++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if l.Front() == i {
		l.front = i.Next
	}

	if l.Back() == i {
		l.back = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.Front() {
		return
	}

	i.Prev.Next = i.Next

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if i == l.Back() {
		l.back = i.Prev
	}

	i.Next = l.Front()
	l.Front().Prev = i
	l.front = i
}

func NewList() *list { //nolint:revive
	return new(list)
}
