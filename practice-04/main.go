package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"example.com/practice-04/internal/task"
	myMW "example.com/practice-04/pkg/middleware"
)

func main() {
	repo := task.NewRepo("tasks.json")
	h := task.NewHandler(repo)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(myMW.Logger)
	r.Use(myMW.SimpleCORS)
	
	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/tasks", h.Routes())
	})
	
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}