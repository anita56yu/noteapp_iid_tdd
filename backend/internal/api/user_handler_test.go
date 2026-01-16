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
	"github.com/golang-jwt/jwt/v5"
)

var testJwtSecret = []byte("supersecretkey")

func setupUserTest() (*chi.Mux, *useruc.UserUsecase, useruc.UserRepository) {
	repo := userrepo.NewInMemoryUserRepository()
	mapper := &useruc.UserMapper{}
	uc := useruc.NewUserUsecase(repo, mapper)
	handler := NewUserHandler(uc, testJwtSecret)

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
	uc.Register("", "existinguser", "password123")

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

func TestUserHandler_Login_ReturnsJWT_Success(t *testing.T) {
	router, uc, _ := setupUserTest()
	registeredUser, _ := uc.Register("", "testuser", "password123") // Register user first

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

	var res LoginResponse
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Token == "" {
		t.Error("expected a non-empty token in response")
	}

	// Validate the JWT token
	token, err := jwt.Parse(res.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return testJwtSecret, nil
	})

	if err != nil {
		t.Fatalf("failed to parse JWT token: %v", err)
	}

	if !token.Valid {
		t.Error("generated token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("failed to get claims from token")
	}

	if claims["user_id"] != registeredUser.ID {
		t.Errorf("expected user_id %s in token claims, got %s", registeredUser.ID, claims["user_id"])
	}

	if res.UserID != registeredUser.ID {
		t.Errorf("expected UserID %s in response, got %s", registeredUser.ID, res.UserID)
	}
	if res.Username != registeredUser.Username {
		t.Errorf("expected Username %s in response, got %s", registeredUser.Username, res.Username)
	}
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	router, uc, _ := setupUserTest()
	uc.Register("", "testuser", "password123")

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
