package converters

import (
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
)

func FlatFromGen(reqGen generated.PostFlatCreateJSONBody) entity.Flat {
	return entity.Flat{
		HouseID: reqGen.HouseId,
		Price:   reqGen.Price,
		Rooms:   *reqGen.Rooms,
	}
}

func FlatToGen(flat entity.Flat) generated.Flat {
	return generated.Flat{
		HouseId: flat.HouseID,
		Id:      flat.ID,
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  generated.Status(flat.Status),
	}
}
