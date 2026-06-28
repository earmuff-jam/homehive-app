package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

// Embed the React production build.
//
//go:embed web/dist/*
var frontend embed.FS

func main() {
	// Get the embedded dist directory.
	dist, err := fs.Sub(frontend, "web/dist")
	if err != nil {
		log.Fatalf("failed to load embedded frontend: %v", err)
	}

	mux := http.NewServeMux()

	// API routes.
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello from Go"))
	})

	// React application.
	mux.Handle("/", spaHandler(dist))

	log.Println("Server listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func spaHandler(files fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(files))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")

		if path == "" {
			path = "index.html"
		}

		if _, err := files.Open(path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		r.URL.Path = "/index.html"
		fileServer.ServeHTTP(w, r)
	})
}
