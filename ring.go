package tchart

// RingBuffer ring buffer
type RingBuffer struct {
	buffer   []interface{}
	length   int
	capacity int
	tail     int
}

// NewRingBuffer new ring buffer
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		buffer:   make([]interface{}, capacity, capacity),
		length:   0,
		capacity: capacity,
		tail:     0,
	}
}

// Len length of ring buffer
func (r *RingBuffer) Len() int {
	return r.length
}

// Capacity capacity of ring buffer
func (r *RingBuffer) Capacity() int {
	return r.capacity
}

// Add add new element
func (r *RingBuffer) Add(v interface{}) {
	if r.length < r.capacity {
		r.length++
	}
	r.buffer[r.tail] = v
	r.tail = (r.tail + 1) % r.capacity
}

// Slice return slice of ring buffer
func (r *RingBuffer) Slice(i, j int) []interface{} {
	if r.length < r.capacity {
		j = Min(j, r.length)
		return r.buffer[i:j]
	}
	s := append(r.buffer[r.tail:r.capacity], r.buffer[:r.tail]...)
	return s[i:j]
}

// Tail return last n elements of ring buffer
func (r *RingBuffer) Tail(n int) []interface{} {
	start := Max(0, r.length-n)

	return r.Slice(start, r.length)
}
