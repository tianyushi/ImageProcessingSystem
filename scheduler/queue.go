// Citation: Dynamic Circular Work-Stealing Dequeue by David Chase & Yossi Lev; Re-Implemented in GO
package scheduler

import (
	"sync/atomic"
)

type Buffer struct {
	tasks      []interface{}
	bufferSize int
}

func CreateBuffer(bufferSize int) *Buffer {
	capacity := 1 << bufferSize
	return &Buffer{
		tasks:      make([]interface{}, capacity),
		bufferSize: bufferSize,
	}
}

func (buffer *Buffer) Capacity() int {
	return 1 << buffer.bufferSize
}

func (buffer *Buffer) Retrieve(index int) interface{} {
	return buffer.tasks[index%buffer.Capacity()]
}

func (buffer *Buffer) Store(index int, task interface{}) {
	buffer.tasks[index%buffer.Capacity()] = task
}

func (buffer *Buffer) Expand(lower, upper int64) *Buffer {
	expandedBuffer := CreateBuffer(buffer.bufferSize + 1)
	for i := upper; i < lower; i++ {
		expandedBuffer.Store(int(i), buffer.Retrieve(int(i)))
	}
	return expandedBuffer
}

type ConcurrentDeque struct {
	bottom int64
	top    int64
	buffer *Buffer
}

func NewDequeue(bufferSize int) *ConcurrentDeque {
	return &ConcurrentDeque{
		buffer: CreateBuffer(bufferSize),
	}
}

func (dq *ConcurrentDeque) PushBottom(item interface{}) {
	bot := atomic.LoadInt64(&dq.bottom)
	top := atomic.LoadInt64(&dq.top)
	buf := dq.buffer
	diff := bot - top
	if diff >= int64(buf.Capacity()-1) {
		buf = buf.Expand(bot, top)
		dq.buffer = buf
	}
	buf.Store(int(bot), item)
	atomic.AddInt64(&dq.bottom, 1)
}

func (dq *ConcurrentDeque) PopBottom() interface{} {
	bot := atomic.LoadInt64(&dq.bottom) - 1
	atomic.StoreInt64(&dq.bottom, bot)
	top := atomic.LoadInt64(&dq.top)
	diff := bot - top
	if diff < 0 {
		atomic.StoreInt64(&dq.bottom, top)
		return nil
	}
	item := dq.buffer.Retrieve(int(bot))
	if diff > 0 {
		if diff < int64(dq.buffer.Capacity())/4 {
			contractedBuffer := dq.buffer.Shrink(bot, top)
			dq.buffer = contractedBuffer
		}
		return item
	}
	if !atomic.CompareAndSwapInt64(&dq.top, top, top+1) {
		item = nil
	}
	atomic.StoreInt64(&dq.bottom, top+1)
	return item
}

func (dq *ConcurrentDeque) Steal() interface{} {
	top := atomic.LoadInt64(&dq.top)
	buf := dq.buffer
	bot := atomic.LoadInt64(&dq.bottom)
	diff := bot - top
	if diff <= 0 {
		return nil
	}
	item := buf.Retrieve(int(top))
	if !atomic.CompareAndSwapInt64(&dq.top, top, top+1) {
		return nil
	}
	return item
}

func (buffer *Buffer) Shrink(lower, upper int64) *Buffer {
	newbufferSize := buffer.bufferSize - 1
	if newbufferSize < 1 {
		newbufferSize = 1
	}
	contractedBuffer := CreateBuffer(newbufferSize)
	for i := upper; i < lower; i++ {
		contractedBuffer.Store(int(i), buffer.Retrieve(int(i)))
	}
	return contractedBuffer
}
