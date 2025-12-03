package main

import (
	"log"
	"net/http"

	"noteapp/internal/api"
	"noteapp/internal/middleware"
	"noteapp/internal/repository/contentrepo"
	"noteapp/internal/repository/noterepo"
	"noteapp/internal/repository/userrepo"
	"noteapp/internal/usecase/contentuc"
	"noteapp/internal/usecase/noteuc"
	"noteapp/internal/usecase/useruc"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// TODO: Move this to a configuration or environment variable in a production environment.
	var jwtSecret = []byte("supersecretkey")

	// 1. Dependency Injection
	noteRepo := noterepo.NewInMemoryNoteRepository()
	contentRepo := contentrepo.NewInMemoryContentRepository()
	userRepo := userrepo.NewInMemoryUserRepository()

	noteUsecase := noteuc.NewNoteUsecase(noteRepo)
	contentUsecase := contentuc.NewContentUsecase(contentRepo)
	userMapper := &useruc.UserMapper{}
	userUsecase := useruc.NewUserUsecase(userRepo, userMapper)

	// Create a test user for development
	_, err := userUsecase.Register("testUser1", "testuser", "password")
	if err != nil {
		log.Fatalf("Failed to create test user: %v", err)
	}

	noteHandler := api.NewNoteHandler(noteUsecase, contentUsecase, userUsecase)
	userHandler := api.NewUserHandler(userUsecase, jwtSecret)

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
	router.Use(chiMiddleware.Logger) // Add a logger middleware

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

	// Public routes
	router.Post("/register", userHandler.Register)
	router.Post("/login", userHandler.Login)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.NewAuthMiddleware(jwtSecret)) // Use our new AuthMiddleware

		r.Post("/notes", noteHandler.CreateNote)
		r.Get("/notes/{id}", noteHandler.GetNoteByID)
		r.Delete("/notes/{id}", noteHandler.DeleteNote)
		r.Put("/notes/{id}", noteHandler.UpdateNote)
		r.Post("/notes/{id}/contents", noteHandler.AddContent)
		r.Put("/notes/{id}/contents/{contentId}", noteHandler.UpdateContent)
		r.Delete("/notes/{id}/contents/{contentId}", noteHandler.DeleteContent)
		r.Get("/notes/accessible-notes", noteHandler.GetAccessibleNotesForUser)
		r.Post("/notes/{noteID}/keywords", noteHandler.TagNote)
		r.Get("/notes/keywords", noteHandler.FindNotesByKeyword)
		r.Delete("/notes/{noteID}/keywords/{keyword}", noteHandler.UntagNote)
		r.Post("/notes/{noteID}/shares", noteHandler.ShareNote)
		r.Delete("/notes/{noteID}/shares", noteHandler.RevokeAccess)
		r.Get("/notes/{noteID}/ws", noteHandler.HandleWebSocket)
	})

	// 3. Server Startup
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
