package useruc

import (
	"noteapp/internal/domain/user"
	"noteapp/internal/repository/userrepo"

	"golang.org/x/crypto/bcrypt"
)

// UserUsecase provides business logic for user management.
type UserUsecase struct {
	repo   userrepo.UserRepository
	mapper *UserMapper
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(repo userrepo.UserRepository, mapper *UserMapper) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		mapper: mapper,
	}
}

// Register registers a new user.
func (uc *UserUsecase) Register(username, password string) (*UserDTO, error) {
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

	domainUser, err := user.NewUser("", username, string(hashedPassword))
	if err != nil {
		return nil, mapDomainError(err)
	}

	userPO := uc.mapper.FromDomain(domainUser)
	err = uc.repo.Save(userPO)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	return &UserDTO{
		ID:       domainUser.ID(),
		Username: domainUser.Username(),
	}, nil
}

// Login authenticates a user.
func (uc *UserUsecase) Login(username, password string) (*UserDTO, error) {
	userPO, err := uc.repo.FindByUsername(username)
	if err != nil {
		return nil, mapRepositoryError(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(userPO.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials // Password mismatch
	}

	domainUser, err := uc.mapper.ToDomain(userPO)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		ID:       domainUser.ID(),
		Username: domainUser.Username(),
	}, nil
}
