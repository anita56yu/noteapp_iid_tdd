package userrepo

import (
	"errors"
	"noteapp/internal/usecase/useruc"
	"sync"
)

// InMemoryUserRepository is an in-memory implementation of UserRepository.
type InMemoryUserRepository struct {
	users map[string]*useruc.UserPO
	mu    sync.RWMutex
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*useruc.UserPO),
	}
}

// Save saves a user to the repository.
func (r *InMemoryUserRepository) Save(user *useruc.UserPO) error {
	if user == nil {
		// You might want to define a specific error for this
		return errors.New("user cannot be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate username
	for _, existingUser := range r.users {
		if existingUser.Username == user.Username && existingUser.ID != user.ID {
			return useruc.ErrUsernameExists
		}
	}

	// Create a deep copy to store
	storedUser := &useruc.UserPO{
		ID:                user.ID,
		Username:          user.Username,
		PasswordHash:      user.PasswordHash,
		AccessibleNoteIDs: make([]string, len(user.AccessibleNoteIDs)),
	}
	copy(storedUser.AccessibleNoteIDs, user.AccessibleNoteIDs)

	r.users[user.ID] = storedUser
	return nil
}

// FindByID retrieves a user by their ID.
func (r *InMemoryUserRepository) FindByID(id string) (*useruc.UserPO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, useruc.ErrInvalidCredentials
	}

	return r.deepCopy(user), nil
}

// FindByUsername retrieves a user by their username.
func (r *InMemoryUserRepository) FindByUsername(username string) (*useruc.UserPO, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return r.deepCopy(user), nil
		}
	}

	return nil, useruc.ErrInvalidCredentials
}

func (r *InMemoryUserRepository) deepCopy(user *useruc.UserPO) *useruc.UserPO {
	newUser := &useruc.UserPO{
		ID:                user.ID,
		Username:          user.Username,
		PasswordHash:      user.PasswordHash,
		AccessibleNoteIDs: make([]string, len(user.AccessibleNoteIDs)),
	}
	copy(newUser.AccessibleNoteIDs, user.AccessibleNoteIDs)
	return newUser
}
