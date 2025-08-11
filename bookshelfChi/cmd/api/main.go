package main

import (
	"bookshelfChi/internal/handlers"
	"bookshelfChi/internal/mw"
	"bookshelfChi/internal/store"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	addr := flag.String("addr", ":8081", "HTTP listen address")
	flag.Parse()

	// 1) ОДИН раз создаём хранилище и хендлеры
	mem := store.NewMemory()
	h := handlers.New(mem, log.Default())

	// 2) Роутер и middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mw.Logging(log.Default()))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// 3) Маршруты (без завершающего слэша у {id})
	r.Route("/v1", func(r chi.Router) {
		r.Get("/books", h.ListBooks)
		r.Post("/books", h.CreateBook)
		r.Get("/books/{id}", h.GetBook)
		r.Put("/books/{id}", h.UpdateBook)
		r.Delete("/books/{id}", h.DeleteBook)
	})

	// 4) Сервер + graceful shutdown
	srv := &http.Server{
		Addr:              *addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", *addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("server stopped")
}
