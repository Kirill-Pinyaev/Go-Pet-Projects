package store

import (
	"bookshelfChi/internal/model"
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("book not found")

type Memory struct {
	mu     sync.RWMutex
	data   map[int64]model.Book
	nextID int64
}

func NewMemory() *Memory {
	return &Memory{
		data:   make(map[int64]model.Book),
		nextID: 1,
	}
}

func (m *Memory) List(offset, limit int) []model.Book {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]model.Book, 0, len(m.data))
	for _, b := range m.data {
		res = append(res, b)
	}

	start := min(offset, len(res))
	end := min(start+limit, len(res))

	return res[start:end]
}

func (m *Memory) Get(id int64) (model.Book, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, ok := m.data[id]

	if !ok {
		return model.Book{}, ErrNotFound
	}

	return b, nil
}

func (m *Memory) Create(dto model.CreateBookDTO) model.Book {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	id := m.nextID
	m.nextID++
	b := model.Book{
		ID:        id,
		Title:     dto.Title,
		Author:    dto.Author,
		Year:      dto.Year,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return b
}

func (m *Memory) Update(id int64, dto model.UpdateBookDTO) (model.Book, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, ok := m.data[id]

	if !ok {
		return model.Book{}, ErrNotFound
	}

	b.Title = dto.Title
	b.Author = dto.Author
	b.Year = dto.Year

	b.UpdatedAt = time.Now()

	m.data[id] = b

	return b, nil

}

func (m *Memory) Delete(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.data[id]

	if !ok {
		return ErrNotFound
	}

	delete(m.data, id)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
