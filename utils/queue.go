package utils

import (
	"runtime"
	"sync"
)

type FixedQueue[T any] struct {
	elements []T
	front    int // 指向队列头
	rear     int // 指向下一个插入位置
	size     int // 当前元素数量
	capacity int // 队列容量

	rwLock *sync.RWMutex
}

// NewFixedQueue 创建一个固定长度的队列
func NewFixedQueue[T any](capacity int) *FixedQueue[T] {
	if capacity <= 0 {
		capacity = 1024
	}
	return &FixedQueue[T]{
		elements: make([]T, capacity),
		capacity: capacity,
		rwLock:   &sync.RWMutex{},
	}
}

// Enqueue 添加元素到队列尾部，当队列满时覆盖最旧的元素
func (q *FixedQueue[T]) Enqueue(item T) {
	q.rwLock.Lock()
	defer q.rwLock.Unlock()
	if q.size < q.capacity {
		q.elements[q.rear] = item
		q.rear = (q.rear + 1) % q.capacity
		q.size++
		return
	}
	// 覆盖最旧的元素
	q.elements[q.rear] = item
	q.rear = (q.rear + 1) % q.capacity
	q.front = (q.front + 1) % q.capacity
}

// Dequeue 移除并返回队列头部的元素
func (q *FixedQueue[T]) Dequeue() interface{} {
	if q.size == 0 {
		return nil
	}
	q.rwLock.Lock()
	defer q.rwLock.Unlock()
	item := q.elements[q.front]
	// q.elements[q.front] = nil // 避免内存泄漏
	q.front = (q.front + 1) % q.capacity
	q.size--
	return item
}

// Peek 返回队列头部的元素但不移除
func (q *FixedQueue[T]) Peek() interface{} {
	if q.size == 0 {
		return nil
	}
	q.rwLock.RLock()
	defer q.rwLock.RUnlock()
	return q.elements[q.front]
}

// Size 返回队列当前元素数量
func (q *FixedQueue[T]) Size() int {
	return q.size
}

// Capacity 返回队列容量
func (q *FixedQueue[T]) Capacity() int {
	return q.capacity
}

// ToSlice 将队列元素转换为切片（按队列顺序）
func (q *FixedQueue[T]) ToSlice(reverse bool) []T {
	q.rwLock.RLock()
	defer q.rwLock.RUnlock()
	slice := make([]T, q.size)
	for i := range q.size {
		pos := (q.front + i) % q.capacity
		if reverse {
			slice[q.size-1-i] = q.elements[pos]
		} else {
			slice[i] = q.elements[pos]
		}
	}
	return slice
}

// Flush 惰性清空
func (q *FixedQueue[T]) Flush(gc bool) {
	q.rwLock.Lock()
	defer q.rwLock.Unlock()
	q.front = 0
	q.rear = 0
	q.size = 0
	if gc {
		q.elements = make([]T, q.capacity)
		runtime.GC()
	}
}
