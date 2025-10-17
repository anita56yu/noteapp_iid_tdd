package main

import (
	"log"
	"net/http"

	"noteapp/internal/api"
	"noteapp/internal/repository"
	"noteapp/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Dependency Injection
	repo := repository.NewInMemoryNoteRepository()
	noteUsecase := usecase.NewNoteUsecase(repo)
	noteHandler := api.NewNoteHandler(noteUsecase)

	// 2. Routing
	router := chi.NewRouter()
	router.Use(middleware.Logger) // Add a logger middleware
	router.Post("/notes", noteHandler.CreateNote)
	router.Get("/notes/{id}", noteHandler.GetNoteByID)
	router.Delete("/notes/{id}", noteHandler.DeleteNote)
	router.Post("/notes/{id}/contents", noteHandler.AddContent)
	router.Put("/notes/{id}/contents/{contentId}", noteHandler.UpdateContent)
	router.Delete("/notes/{id}/contents/{contentId}", noteHandler.DeleteContent)
	router.Get("/users/{userID}/accessible-notes", noteHandler.GetAccessibleNotesForUser)
	router.Delete("/users/{ownerID}/notes/{noteID}/shares", noteHandler.RevokeAccess)

	// 3. Server Startup
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
