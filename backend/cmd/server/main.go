package main

import (
	"log"
	"net/http"

	"noteapp/internal/api"
	"noteapp/internal/repository/contentrepo"
	"noteapp/internal/repository/noterepo"
	"noteapp/internal/usecase/contentuc"
	"noteapp/internal/usecase/noteuc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// 1. Dependency Injection
	noteRepo := noterepo.NewInMemoryNoteRepository()
	contentRepo := contentrepo.NewInMemoryContentRepository()

	noteUsecase := noteuc.NewNoteUsecase(noteRepo)
	contentUsecase := contentuc.NewContentUsecase(contentRepo)

	noteHandler := api.NewNoteHandler(noteUsecase, contentUsecase)

	// test data
	n1, err := noteUsecase.CreateNote("", "Test Note 1", "testUser1")
	if err != nil {
		log.Fatalf("Failed to create test note: %v", err)
	}
	_, err = noteUsecase.CreateNote("", "Test Note 2", "testUser1")
	if err != nil {
		log.Fatalf("Failed to create test note: %v", err)
	}
	c1, err := contentUsecase.CreateContent(n1, "", "Content for Note 1", "text")
	if err != nil {
		log.Fatalf("Failed to create test content: %v", err)
	}
	noteUsecase.AddContent(n1, c1, -1, 0)

	// 2. Routing
	router := chi.NewRouter()
	router.Use(middleware.Logger) // Add a logger middleware

	// Custom middleware to log the Origin header
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				log.Printf("Request Origin: %s", origin)
			}
			next.ServeHTTP(w, r)
		})
	})

	// Add CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "vscode-file://vscode-app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	}))

	router.Post("/notes", noteHandler.CreateNote)
	router.Get("/notes/{id}", noteHandler.GetNoteByID)
	router.Delete("/notes/{id}", noteHandler.DeleteNote)
	router.Put("/notes/{id}", noteHandler.UpdateNote)
	router.Post("/notes/{id}/contents", noteHandler.AddContent)
	router.Put("/notes/{id}/contents/{contentId}", noteHandler.UpdateContent)
	router.Delete("/notes/{id}/contents/{contentId}", noteHandler.DeleteContent)
	router.Get("/users/{userID}/accessible-notes", noteHandler.GetAccessibleNotesForUser)
	router.Delete("/users/{ownerID}/notes/{noteID}/shares", noteHandler.RevokeAccess)
	router.Get("/notes/{noteID}/ws", noteHandler.HandleWebSocket)

	// 3. Server Startup
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
