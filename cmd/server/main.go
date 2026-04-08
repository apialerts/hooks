package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apialerts/hooks/internal/cleanup"
	"github.com/apialerts/hooks/internal/db"
	"github.com/apialerts/hooks/internal/handler"
	"github.com/apialerts/hooks/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFS embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://hooks:hooks@localhost:5432/hooks?sslmode=disable"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%s", port)
	}

	ctx := context.Background()

	pool, err := db.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(databaseURL); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	queries := db.New(pool)

	staticContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("failed to setup static fs: %v", err)
	}

	h := handler.New(queries, baseURL)
	rl := middleware.NewRateLimiter()

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticContent))))

	r.Get("/", h.Home)
	r.Post("/endpoints", h.CreateEndpoint)
	r.Get("/privacy", h.Privacy)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r.Get("/robots.txt", h.RobotsTxt)
	r.Get("/sitemap.xml", h.SitemapXml)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(rl.Limit)
		r.Get("/requests", h.PollRequests)
		r.Get("/requests/{reqId}", h.RequestDetail)
		r.Put("/config", h.UpdateConfig)
		r.Delete("/delete", h.DeleteEndpoint)
		r.Post("/test", h.TestEndpoint)
		r.HandleFunc("/*", h.HandleEndpoint)
		r.HandleFunc("/", h.HandleEndpoint)
	})

	cleanupJob := cleanup.New(queries)
	cleanupJob.Start(ctx)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("starting server on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
