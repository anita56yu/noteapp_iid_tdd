package useruc

import (
	"noteapp/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
)

// UserUsecase provides business logic for user management.
type UserUsecase struct {
	repo   UserRepository
	mapper *UserMapper
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(repo UserRepository, mapper *UserMapper) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		mapper: mapper,
	}
}

// Register registers a new user.
func (uc *UserUsecase) Register(id, username, password string) (*UserDTO, error) {
	_, err := uc.repo.FindByUsername(username)
	if err == nil {
		return nil, ErrUsernameExists
	}
	if password == "" {
		return nil, ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err // Internal error during password hashing
	}

	domainUser, err := user.NewUser(id, username, string(hashedPassword))
	if err != nil {
		return nil, mapDomainError(err)
	}

	userPO := uc.mapper.FromDomain(domainUser)
	err = uc.repo.Save(userPO)
	if err != nil {
		return nil, err
	}

	return uc.mapper.ToDTO(domainUser), nil
}

// Login authenticates a user.
func (uc *UserUsecase) Login(username, password string) (*UserDTO, error) {
	userPO, err := uc.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	domainUser, err := uc.mapper.ToDomain(userPO)
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(domainUser.PasswordHash()), []byte(password)) != nil {
		return nil, ErrInvalidCredentials // Password mismatch
	}

	return uc.mapper.ToDTO(domainUser), nil
}

// AddAccessibleNote adds a note ID to the user's accessible notes.
func (uc *UserUsecase) AddAccessibleNote(userID, noteID string) error {
	userPO, err := uc.repo.FindByID(userID)
	if err != nil {
		return err
	}

	domainUser, err := uc.mapper.ToDomain(userPO)
	if err != nil {
		return err
	}

	domainUser.AddAccessibleNoteID(noteID)

	updatedUserPO := uc.mapper.FromDomain(domainUser)
	err = uc.repo.Save(updatedUserPO)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAccessibleNote removes a note ID from the user's accessible notes.
func (uc *UserUsecase) RemoveAccessibleNote(userID, noteID string) error {
	userPO, err := uc.repo.FindByID(userID)
	if err != nil {
		return err
	}

	domainUser, err := uc.mapper.ToDomain(userPO)
	if err != nil {
		return err
	}

	domainUser.RemoveAccessibleNoteID(noteID)

	updatedUserPO := uc.mapper.FromDomain(domainUser)
	err = uc.repo.Save(updatedUserPO)
	if err != nil {
		return err
	}

	return nil
}

// CheckUser checks if a user exists by their ID.
func (uc *UserUsecase) CheckUser(userID string) (bool, error) {
	_, err := uc.repo.FindByID(userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (uc *UserUsecase) FindUserIDByUsername(username string) (string, error) {
	userPO, err := uc.repo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	return userPO.ID, nil
}
