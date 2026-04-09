package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
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
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbHost := os.Getenv("DB_HOST")
		dbName := os.Getenv("DB_NAME")
		dbPort := "5432"
		dbPortTemp := os.Getenv("DB_PORT")
		if dbPortTemp != "" {
			dbPort = dbPortTemp
		}
		dbSSLMode := os.Getenv("DB_SSLMODE")
		if dbSSLMode == "" {
			dbSSLMode = "require"
		}
		if dbUser != "" && dbPassword != "" && dbHost != "" && dbName != "" {
			databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", url.QueryEscape(dbUser), url.QueryEscape(dbPassword), dbHost, dbPort, dbName, dbSSLMode)
		} else {
			databaseURL = fmt.Sprintf("postgres://hooks:hooks@localhost:%s/hooks?sslmode=disable", dbPort)
		}
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%s", port)
	}

	killSwitch := os.Getenv("KILL_SWITCH") == "true"

	ctx := context.Background()

	staticContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("failed to setup static fs: %v", err)
	}

	// Kill switch — serve a maintenance page without connecting to the database
	if killSwitch {
		log.Println("KILL_SWITCH enabled — serving maintenance page only")

		r := chi.NewRouter()
		r.Use(chimw.Logger)
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticContent))))
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("maintenance"))
		})
		r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Maintenance - hooks.apialerts.com</title>
    <link rel="stylesheet" href="/static/styles.css">
</head>
<body class="bg-gray-50 dark:bg-dark-bg min-h-screen flex items-center justify-center font-sans text-gray-900 dark:text-dark-text">
    <div class="text-center px-4">
        <h1 class="text-2xl font-bold mb-3">Temporarily Unavailable</h1>
        <p class="text-gray-600 dark:text-dark-text-secondary mb-6">hooks.apialerts.com is undergoing maintenance and will be back shortly.</p>
        <a href="https://apialerts.com" class="text-brand hover:underline text-sm font-medium">apialerts.com</a>
    </div>
</body>
</html>`))
		})

		srv := &http.Server{Addr: ":" + port, Handler: r}
		go func() {
			log.Printf("maintenance server on :%s", port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("server error: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		srv.Shutdown(context.Background())
		return
	}

	pool, err := db.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(databaseURL); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	queries := db.New(pool)

	h := handler.New(queries, baseURL)
	rl := middleware.NewRateLimiter()

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.CORS)

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
