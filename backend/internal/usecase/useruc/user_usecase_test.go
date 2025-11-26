package useruc_test

import (
	"errors"
	"noteapp/internal/repository/userrepo"
	"noteapp/internal/usecase/useruc"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase_Register_Success(t *testing.T) {
	// Arrange
	repo := userrepo.NewInMemoryUserRepository()
	mapper := &useruc.UserMapper{}
	usecase := useruc.NewUserUsecase(repo, mapper)

	username := "newuser"
	password := "password123"

	// Act
	userDTO, err := usecase.Register(username, password)

	// Assert
	if err != nil {
		t.Fatalf("Register() returned an unexpected error: %v", err)
	}
	if userDTO == nil {
		t.Fatal("Register() returned nil UserDTO")
	}
	if userDTO.Username != username {
		t.Errorf("Expected username %s, got %s", username, userDTO.Username)
	}
	if userDTO.ID == "" {
		t.Error("Expected a non-empty user ID")
	}

	// Verify user saved in repository
	foundPO, err := repo.FindByID(userDTO.ID)
	if err != nil {
		t.Fatalf("FindByID() returned an error: %v", err)
	}
	if foundPO.Username != username {
		t.Errorf("Expected saved user username %s, got %s", username, foundPO.Username)
	}
	if bcrypt.CompareHashAndPassword([]byte(foundPO.PasswordHash), []byte(password)) != nil {
		t.Error("Password hash mismatch for saved user")
	}
}

func TestUserUsecase_Register_UsernameExists(t *testing.T) {
	// Arrange
	repo := userrepo.NewInMemoryUserRepository()
	mapper := &useruc.UserMapper{}
	usecase := useruc.NewUserUsecase(repo, mapper)

	username := "existinguser"
	password := "password123"
	// Register the first user successfully
	_, err := usecase.Register(username, password)
	if err != nil {
		t.Fatalf("Initial register failed: %v", err)
	}

	// Act: Try to register with the same username again
	userDTO, err := usecase.Register(username, password)

	// Assert
	if !errors.Is(err, useruc.ErrUsernameExists) {
		t.Errorf("Expected error %v, got %v", useruc.ErrUsernameExists, err)
	}
	if userDTO != nil {
		t.Error("Expected nil UserDTO, got non-nil")
	}
}

func TestUserUsecase_Register_InvalidInput(t *testing.T) {
	// Arrange
	repo := userrepo.NewInMemoryUserRepository()
	mapper := &useruc.UserMapper{}
	usecase := useruc.NewUserUsecase(repo, mapper)

	// Test empty username
	userDTO, err := usecase.Register("", "password123")
	if err == nil || !errors.Is(err, useruc.ErrInvalidCredentials) {
		t.Errorf("Expected error %v for empty username, got: %v", useruc.ErrInvalidCredentials, err)
	}
	if userDTO != nil {
		t.Error("Expected nil UserDTO for empty username")
	}

	// Test empty password
	userDTO, err = usecase.Register("testuser", "")
	if err == nil || !errors.Is(err, useruc.ErrInvalidCredentials) {
		t.Errorf("Expected error %v for empty password, got: %v", useruc.ErrInvalidCredentials, err)
	}
	if userDTO != nil {
		t.Error("Expected nil UserDTO for empty password")
	}
}

func TestUserUsecase_Login_Success(t *testing.T) {
	// Arrange
	repo := userrepo.NewInMemoryUserRepository()
	mapper := &useruc.UserMapper{}
	usecase := useruc.NewUserUsecase(repo, mapper)

	username := "testuser"
	password := "password123"

	// Register the user first
	registeredUser, err := usecase.Register(username, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Act
	userDTO, err := usecase.Login(username, password)

	// Assert
	if err != nil {
		t.Fatalf("Login() returned an unexpected error: %v", err)
	}
	if userDTO == nil {
		t.Fatal("Login() returned nil UserDTO")
	}
	if userDTO.Username != username {
		t.Errorf("Expected username %s, got %s", username, userDTO.Username)
	}
	if userDTO.ID != registeredUser.ID {
		t.Errorf("Expected ID %s, got %s", registeredUser.ID, userDTO.ID)
	}
}

func TestUserUsecase_Login_InvalidUsername(t *testing.T) {
	// Arrange
	repo := userrepo.NewInMemoryUserRepository()
	mapper := &useruc.UserMapper{}
	usecase := useruc.NewUserUsecase(repo, mapper)

	username := "nonexistent"
	password := "password123"

	// Act
	userDTO, err := usecase.Login(username, password)

	// Assert
	if !errors.Is(err, useruc.ErrInvalidCredentials) {
		t.Errorf("Expected error %v, got %v", useruc.ErrInvalidCredentials, err)
	}
	if userDTO != nil {
		t.Error("Expected nil UserDTO")
	}
}

func TestUserUsecase_Login_InvalidPassword(t *testing.T) {
	// Arrange
	repo := userrepo.NewInMemoryUserRepository()
	mapper := &useruc.UserMapper{}
	usecase := useruc.NewUserUsecase(repo, mapper)

	username := "testuser"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"

	// Register the user first
	_, err := usecase.Register(username, correctPassword)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Act
	userDTO, err := usecase.Login(username, wrongPassword)

	// Assert
	if !errors.Is(err, useruc.ErrInvalidCredentials) {
		t.Errorf("Expected error %v, got %v", useruc.ErrInvalidCredentials, err)
	}
	if userDTO != nil {
		t.Error("Expected nil UserDTO")
	}
}
