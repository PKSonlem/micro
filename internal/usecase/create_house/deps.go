package create_house

import "github.com/timurzdev/mentorship-test-task/internal/entity"

// go: generate mockgen -source=deps.go -destination=mock/deps.go -package=mock
type Repository interface {
	CreateHouse(flat entity.House) error
}
