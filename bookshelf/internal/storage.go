package internal

import (
	"errors"
	"sync"
)

type Book struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author,omitempty"`
}

type Store struct {
	mu    sync.RWMutex
	seq   int64
	books []Book
}

func NewSore() *Store {
	return &Store{}
}

func (s *Store) List() []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Book, len(s.books))

	copy(out, s.books)
	return out
}

func (s *Store) Create(book Book) (Book, error) {
	if book.Title == "" {
		return Book{}, errors.New("title required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++

	book.ID = s.seq
	s.books = append(s.books, book)
	return book, nil
}
