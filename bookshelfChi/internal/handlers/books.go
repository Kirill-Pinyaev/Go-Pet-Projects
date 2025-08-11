package handlers

import (
	"bookshelfChi/internal/model"
	"bookshelfChi/internal/store"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) dbg(label string) {
	if m, ok := h.store.(*store.Memory); ok {
		log.Printf("[DBG] %s store_ptr=%p", label, m)
	} else {
		log.Printf("[DBG] %s store_type=%T (не Memory)", label, h.store)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg})
}

type Store interface {
	List(offset, limit int) []model.Book
	Get(id int64) (model.Book, error)
	Create(dto model.CreateBookDTO) model.Book
	Update(id int64, dto model.UpdateBookDTO) (model.Book, error)
	Delete(id int64) error
}

type Logger interface {
	Printf(format string, v ...any)
}

type Handlers struct {
	store Store
	log   Logger
}

func New(s Store, l Logger) *Handlers {
	return &Handlers{
		store: s,
		log:   l,
	}
}

func (h *Handlers) ListBooks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	offset, _ := strconv.Atoi(q.Get("offset"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	if limit <= 0 || limit > 10000 {
		limit = 100
	}

	items := h.store.List(offset, limit)
	writeJSON(w, http.StatusOK, map[string]any{
		"items": items, "offset": offset, "limit": limit,
	})

}

func (h *Handlers) GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	b, err := h.store.Get(id)

	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, store.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *Handlers) CreateBook(w http.ResponseWriter, r *http.Request) {
	var dto model.CreateBookDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "malformed json")
		return
	}
	if err := dto.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	b := h.store.Create(dto)
	w.Header().Set("Location", "/v1/books/"+strconv.FormatInt(b.ID, 10))
	writeJSON(w, http.StatusCreated, b)
}

func (h *Handlers) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var dto model.UpdateBookDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "malformed json")
		return
	}
	if err := dto.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	b, err := h.store.Update(id, dto)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, store.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, b)
}

func (h *Handlers) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.store.Delete(id); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, store.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
