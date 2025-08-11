package jobqueue

import "sync"

type Queue struct {
	mu     sync.Mutex
	cond   *sync.Cond
	data   []string
	closed bool
}

func New() *Queue {
	q := &Queue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Push добавляет URL и пробуждает один воркер.
func (q *Queue) Push(url string) {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	q.data = append(q.data, url)
	q.mu.Unlock()
	q.cond.Signal()
}

// Pop блокируется, пока очередь пуста. Второй bool=false, если очередь закрыта.
func (q *Queue) Pop() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.data) == 0 && !q.closed {
		q.cond.Wait()
	}
	if len(q.data) == 0 && q.closed {
		return "", false
	}
	url := q.data[0]
	q.data = q.data[1:]
	return url, true
}

// Close завершает очередь и будит всех ожидающих.
func (q *Queue) Close() {
	q.mu.Lock()
	q.closed = true
	q.mu.Unlock()
	q.cond.Broadcast()
}
