// Package main is the ramiro.me web server entry point.
package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	ramirome "github.com/rejkpp/ramiro.me"
	"github.com/rejkpp/ramiro.me/internal/handler"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	staticRoot, err := fs.Sub(ramirome.StaticFS, "static")
	if err != nil {
		log.Fatalf("static fs: %v", err)
	}
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticRoot))))

	pub := handler.NewPublic()
	r.Get("/", pub.Home)
	r.Get("/about", pub.About)
	r.Get("/contact", pub.Contact)
	r.Get("/projects", pub.Projects)

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("ramiro.me listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}
