package useruc

import (
	"noteapp/internal/domain/user"
	"noteapp/internal/repository/userrepo"
)

// UserMapper provides methods to convert between domain.User and userrepo.UserPO.
type UserMapper struct{}

// ToDomain converts a userrepo.UserPO to a domain.User.
func (m *UserMapper) ToDomain(po *userrepo.UserPO) (*user.User, error) {
	u, err := user.NewUser(po.ID, po.Username, po.PasswordHash)
	if err != nil {
		return nil, err
	}
	for _, noteID := range po.AccessibleNoteIDs {
		u.AddAccessibleNoteID(noteID)
	}
	return u, nil
}

// FromDomain creates a userrepo.UserPO from a domain.User.
func (m *UserMapper) FromDomain(u *user.User) *userrepo.UserPO {
	return &userrepo.UserPO{
		ID:                u.ID(),
		Username:          u.Username(),
		PasswordHash:      u.PasswordHash(),
		AccessibleNoteIDs: u.AccessibleNoteIDs(),
	}
}

func (m *UserMapper) ToDTO(u *user.User) *UserDTO {
	return &UserDTO{
		ID:                u.ID(),
		Username:          u.Username(),
		AccessibleNoteIDs: u.AccessibleNoteIDs(),
	}
}
