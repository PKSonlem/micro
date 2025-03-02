package entity

type Role string

const (
	admin Role = "Admin"
)

func (r Role) IsAdmin() bool {
	return r == admin
}
