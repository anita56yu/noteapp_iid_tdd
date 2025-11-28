package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var testJwtSecret = []byte("testsecret")

func generateTestToken(userID string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	userID := "test-user-id"
	token, err := generateTestToken(userID, testJwtSecret)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	req, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	handled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled = true
		ctxUserID, ok := GetUserIDFromContext(r.Context())
		if !ok || ctxUserID != userID {
			t.Errorf("Expected user ID '%s' in context, but got '%s' (ok: %v)", userID, ctxUserID, ok)
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := NewAuthMiddleware(testJwtSecret)
	handler := middleware(nextHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	if !handled {
		t.Error("Expected next handler to be called, but it was not")
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled = true
	})

	middleware := NewAuthMiddleware(testJwtSecret)
	handler := middleware(nextHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized, got %d", rr.Code)
	}
	if handled {
		t.Error("Expected next handler not to be called, but it was")
	}
}

func TestAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	req, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "InvalidToken")

	rr := httptest.NewRecorder()

	handled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled = true
	})

	middleware := NewAuthMiddleware(testJwtSecret)
	handler := middleware(nextHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized, got %d", rr.Code)
	}
	if handled {
		t.Error("Expected next handler not to be called, but it was")
	}
}

func TestAuthMiddleware_InvalidSignature(t *testing.T) {
	userID := "test-user-id"
	badSecret := []byte("wrongsecret")
	badTokenString, err := generateTestToken(userID, badSecret)
	if err != nil {
		t.Fatalf("Failed to generate bad test token: %v", err)
	}

	req, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+badTokenString)

	rr := httptest.NewRecorder()

	handled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled = true
	})

	middleware := NewAuthMiddleware(testJwtSecret)
	handler := middleware(nextHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized, got %d", rr.Code)
	}
	if handled {
		t.Error("Expected next handler not to be called, but it was")
	}
}

func TestAuthMiddleware_TokenWithoutUserIDClaim(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(testJwtSecret)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	req, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()

	handled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled = true
	})

	middleware := NewAuthMiddleware(testJwtSecret)
	handler := middleware(nextHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized, got %d", rr.Code)
	}
	if handled {
		t.Error("Expected next handler not to be called, but it was")
	}
}
