package entity

type Role string

const (
	RoleClient    Role = "client"
	RoleModerator Role = "moderator"
)

func (r Role) IsModerator() bool {
	return r == RoleModerator
}

func (r Role) IsClient() bool {
	return r == RoleClient
}

// IsValid проверяет валидность роли
func (r Role) IsValid() bool {
	return r == RoleClient || r == RoleModerator
}
