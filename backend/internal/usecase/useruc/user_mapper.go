package useruc

import (
	"noteapp/internal/domain/user"
)

// UserMapper provides methods to convert between domain.User and useruc.UserPO.
type UserMapper struct{}

// ToDomain converts a useruc.UserPO to a domain.User.
func (m *UserMapper) ToDomain(po *UserPO) (*user.User, error) {
	u, err := user.NewUser(po.ID, po.Username, po.PasswordHash)
	if err != nil {
		return nil, err
	}
	for _, noteID := range po.AccessibleNoteIDs {
		u.AddAccessibleNoteID(noteID)
	}
	return u, nil
}

// FromDomain creates a useruc.UserPO from a domain.User.
func (m *UserMapper) FromDomain(u *user.User) *UserPO {
	return &UserPO{
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
