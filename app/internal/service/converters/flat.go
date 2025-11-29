package converters

import (
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
)

func FlatFromGenCreate(reqGen generated.PostFlatCreateJSONBody) entity.Flat {
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

func FlatFromGenUpdate(renGen generated.PostFlatUpdateJSONBody) entity.Flat {
	status := ""
	if renGen.Status != nil {
		status = string(*renGen.Status)
	}
	return entity.Flat{
		ID:     renGen.Id,
		Status: status,
	}
}

func FlatToGenArr(flats []entity.Flat) []generated.Flat {
	result := make([]generated.Flat, 0, len(flats))
	for _, flat := range flats {
		result = append(result, generated.Flat{
			HouseId: flat.HouseID,
			Id:      flat.ID,
			Price:   flat.Price,
			Rooms:   flat.Rooms,
			Status:  generated.Status(flat.Status),
		})
	}
	return result
}
