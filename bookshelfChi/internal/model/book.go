package model

import (
	"errors"
	"strings"
	"time"
)

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBookDTO struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

type UpdateBookDTO struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

func (d CreateBookDTO) Validate() error {
	if strings.TrimSpace(d.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(d.Author) == "" {
		return errors.New("author is required")
	}
	if d.Year < 1450 || d.Year > time.Now().Year()+1 {
		return errors.New("year is out of range")
	}
	return nil
}

func (d UpdateBookDTO) Validate() error {
	return CreateBookDTO(d).Validate()
}
