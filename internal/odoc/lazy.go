package odoc

import "sync"

type lazy[T any] struct {
	once  sync.Once
	value T
}

func (l *lazy[T]) Get(fn func() T) T {
	l.once.Do(func() { l.value = fn() })
	return l.value
}
