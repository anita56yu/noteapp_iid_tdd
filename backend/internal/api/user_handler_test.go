package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"noteapp/internal/repository/userrepo"
	"noteapp/internal/usecase/useruc"
	"testing"

	"github.com/go-chi/chi/v5"
)

func setupUserTest() (*chi.Mux, *useruc.UserUsecase, userrepo.UserRepository) {
	repo := userrepo.NewInMemoryUserRepository()
	mapper := &useruc.UserMapper{}
	uc := useruc.NewUserUsecase(repo, mapper)
	handler := NewUserHandler(uc)

	router := chi.NewRouter()
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)

	return router, uc, repo
}

func TestUserHandler_Register_Success(t *testing.T) {
	router, _, repo := setupUserTest()

	reqBody := RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d; got %d", http.StatusCreated, rr.Code)
	}

	var user useruc.UserDTO
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if user.Username != reqBody.Username {
		t.Errorf("expected username %s; got %s", reqBody.Username, user.Username)
	}
	if user.ID == "" {
		t.Error("expected non-empty user ID")
	}

	// Verify user is stored in repository
	storedUser, err := repo.FindByUsername("testuser")
	if err != nil {
		t.Fatalf("failed to find user in repository: %v", err)
	}
	if storedUser.Username != reqBody.Username {
		t.Errorf("expected stored username %s; got %s", reqBody.Username, storedUser.Username)
	}
}

func TestUserHandler_Register_UsernameExists(t *testing.T) {
	router, uc, _ := setupUserTest()
	uc.Register("existinguser", "password123")

	reqBody := RegisterRequest{
		Username: "existinguser",
		Password: "newpassword",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected status %d; got %d", http.StatusConflict, rr.Code)
	}
}

func TestUserHandler_Register_InvalidInput(t *testing.T) {
	router, _, _ := setupUserTest()

	// Empty password
	reqBody := RegisterRequest{
		Username: "testuser",
		Password: "",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestUserHandler_Login_Success(t *testing.T) {
	router, uc, _ := setupUserTest()
	uc.Register("testuser", "password123")

	reqBody := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
	}

	var user useruc.UserDTO
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if user.Username != reqBody.Username {
		t.Errorf("expected username %s; got %s", reqBody.Username, user.Username)
	}
	if user.ID == "" {
		t.Error("expected non-empty user ID")
	}
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	router, uc, _ := setupUserTest()
	uc.Register("testuser", "password123")

	reqBody := LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d; got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestUserHandler_Login_UserNotFound(t *testing.T) {
	router, _, _ := setupUserTest()

	reqBody := LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d; got %d", http.StatusUnauthorized, rr.Code)
	}
}
