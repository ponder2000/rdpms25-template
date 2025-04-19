package collection

import (
	"sync"

	"github.com/ponder2000/rdpms25-template/pkg/util/generic"
)

// Node represents an element in the doubly linked list.
type Node[T any] struct {
	value T
	prev  *Node[T]
	next  *Node[T]
}

// Deque represents a double-ended queue (deque) using a doubly linked list.
type Deque[T any] struct {
	front *Node[T]
	back  *Node[T]
	size  int
	mutex sync.RWMutex // Mutex to ensure thread-safety
}

// NewDeque creates and returns a new empty deque.
func NewDeque[T any]() *Deque[T] {
	return &Deque[T]{}
}

// PushFront adds an element to the front of the deque.
func (d *Deque[T]) PushFront(value T) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	newNode := &Node[T]{value: value}
	if d.front == nil {
		d.front = newNode
		d.back = newNode
	} else {
		newNode.next = d.front
		d.front.prev = newNode
		d.front = newNode
	}
	d.size++
}

// PushBack adds an element to the back of the deque.
func (d *Deque[T]) PushBack(value T) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	newNode := &Node[T]{value: value}
	if d.back == nil {
		d.front = newNode
		d.back = newNode
	} else {
		newNode.prev = d.back
		d.back.next = newNode
		d.back = newNode
	}
	d.size++
}

// PushBackMultiple adds elements to the back of the deque.
func (d *Deque[T]) PushBackMultiple(arr []T) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	newNodes := generic.Mapper(arr, func(v T) *Node[T] { return &Node[T]{value: v} })

	for i := range newNodes {
		nn := newNodes[i]
		if d.back == nil {
			d.front = nn
			d.back = nn
		} else {
			nn.prev = d.back
			d.back.next = nn
			d.back = nn
		}
		d.size++
	}
}

// PopFront removes and returns the element at the front of the deque.
func (d *Deque[T]) PopFront() T {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.size == 0 {
		var zero T
		return zero // Return the zero value for type T
	}
	value := d.front.value
	d.front = d.front.next
	if d.front != nil {
		d.front.prev = nil
	} else {
		d.back = nil // If the deque becomes empty, reset the back pointer
	}
	d.size--
	return value
}

// PopBack removes and returns the element at the back of the deque.
func (d *Deque[T]) PopBack() T {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.size == 0 {
		var zero T
		return zero // Return the zero value for type T
	}
	value := d.back.value
	d.back = d.back.prev
	if d.back != nil {
		d.back.next = nil
	} else {
		d.front = nil // If the deque becomes empty, reset the front pointer
	}
	d.size--
	return value
}

// PeekFront returns the element at the front of the deque without removing it.
func (d *Deque[T]) PeekFront() T {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.size == 0 {
		var zero T
		return zero // Return the zero value for type T
	}
	return d.front.value
}

// PeekBack returns the element at the back of the deque without removing it.
func (d *Deque[T]) PeekBack() T {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.size == 0 {
		var zero T
		return zero // Return the zero value for type T
	}
	return d.back.value
}

// Size returns the number of elements in the deque.
func (d *Deque[T]) Size() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.size
}

func (d *Deque[T]) IsEmpty() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.size == 0
}

func (d *Deque[T]) Clear() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.front = nil
	d.back = nil
	d.size = 0
}
