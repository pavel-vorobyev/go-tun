package util

import "sync"

type Queue[T any] struct {
	data []*T
}

func (q *Queue[T]) Put(i *T) {
	q.data = append(q.data, i)
}

func (q *Queue[T]) Pop() *T {
	if len(q.data) > 1 {
		item := q.data[0]
		q.data = q.data[1:]
		return item
	}
	return nil
}

func (q *Queue[T]) Length() int {
	return len(q.data)
}

type ConcurrentQueue[T any] struct {
	queue Queue[T]
	mutex sync.RWMutex
}

func (q *ConcurrentQueue[T]) Put(i *T) {
	q.mutex.Lock()
	q.queue.Put(i)
	q.mutex.Unlock()
}

func (q *ConcurrentQueue[T]) Pop() *T {
	q.mutex.Lock()
	item := q.queue.Pop()
	q.mutex.Unlock()
	return item
}

func (q *ConcurrentQueue[T]) Length() int {
	q.mutex.RLock()
	length := q.Length()
	q.mutex.RUnlock()
	return length
}
