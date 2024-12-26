package converters

import (
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
)

func HouseFromGen(genReq generated.PostHouseCreateJSONBody) entity.House {
	return entity.House{
		Address:   genReq.Address,
		Year:      genReq.Year,
		Developer: genReq.Developer,
	}
}
