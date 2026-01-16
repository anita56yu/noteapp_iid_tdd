package useruc

// UserRepository defines the interface for user data persistence.
type UserRepository interface {
	Save(user *UserPO) error
	FindByID(id string) (*UserPO, error)
	FindByUsername(username string) (*UserPO, error)
}
