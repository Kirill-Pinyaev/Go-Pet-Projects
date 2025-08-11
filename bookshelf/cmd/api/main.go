package main

import (
	"bookshelf/internal"
	"bookshelf/internal/middleware"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type app struct {
	store  *internal.Store
	logger *log.Logger
}

func main() {
	port := getenv("PORT", "8080")
	addr := ":" + port

	a := &app{
		store:  internal.NewSore(),
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", a.handlePing)
	mux.HandleFunc("GET /books", a.handleBoolsList)
	mux.HandleFunc("POST /books", a.handleBooksCreate)

	handler := middleware.Recover(a.logger)(middleware.Logger(a.logger)(mux))

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	a.logger.Printf("listening on %s", addr)

	ctx, stop := signalContext()
	defer stop()

	errch := make(chan error, 1)
	go func() {
		errch <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(shCtx)
		a.logger.Println("server: shutdown complete")
	case err := <-errch:
		if !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatalf("server: %v", err)
		}
	}
}

func (a *app) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("pong"))
}

func (a *app) handleBoolsList(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, a.store.List())
}

func (a *app) handleBooksCreate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var in internal.Book
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	book, err := a.store.Create(in)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Location", "/books/"+strconv.FormatInt(book.ID, 10))
	writeJson(w, http.StatusCreated, book)
}

func writeJson(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	writeJson(w, code, map[string]string{"error": msg})
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func signalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt)
}
