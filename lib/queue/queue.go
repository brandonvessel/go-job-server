package queue

import (
	"container/list"
	"errors"
	"fmt"
)

type Queue struct {
	queue *list.List
}

// NewQueue creates a new Queue instance
func NewQueue() *Queue {
	q := new(Queue)
	q.queue = list.New()
	fmt.Println("queue length", q.queue.Len())
	return q
}

// Enqueue adds an element to the back of the queue.
func (c *Queue) Enqueue(value string) {
	// add value to queue
	c.queue.PushBack(value)
}

// Dequeue removes the front element of the queue and returns it.
func (c *Queue) Dequeue() (interface{}, error) {
	// if queue is empty, return error
	if c.queue.Len() > 0 {
		// get value from front of queue
		ele := c.queue.Front()
		// remove value from queue
		c.queue.Remove(ele)
		// return value
		return ele.Value, nil
	}
	return nil, errors.New("peep error: queue is empty")
}

// Front returns the value of the front element of the queue.
func (c *Queue) Front() (string, error) {
	// if queue is empty, return error
	if c.queue.Len() > 0 {
		// get value from front of queue
		if val, ok := c.queue.Front().Value.(string); ok {
			// return value
			return val, nil
		}
		return "", errors.New("peep Error: queue datatype is incorrect")
	}
	return "", errors.New("peep error: queue is empty")
}

// Size returns the number of elements in the queue.
func (c *Queue) Size() int {
	fmt.Println("queue size:", (*c.queue).Len())
	return c.queue.Len()
}

// New creates a new queue.
func (c *Queue) Empty() bool {
	return c.Size() == 0
}
