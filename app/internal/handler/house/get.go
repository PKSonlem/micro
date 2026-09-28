package house

import (
	"context"

	"github.com/PKSonlem/micro/internal/generated"
	"github.com/PKSonlem/micro/internal/service/converters"
)

func (h *Handler) GetHouseFlats(ctx context.Context, request generated.GetHouseIdRequestObject) (generated.GetHouseIdResponseObject, error) {
	role, err := h.roleProvider.GetRole(ctx)
	if err != nil {
		h.logger.Error(ctx, err)
		return generated.GetHouseId401Response{}, nil
	}

	isModerator := role.IsModerator()

	flats, err := h.usecase.HandleGetHouseFlatID(ctx, request.Id, isModerator)
	if err != nil {
		h.logger.Error(ctx, err)
		return generated.GetHouseId500JSONResponse{
			N5xxJSONResponse: generated.N5xxJSONResponse{
				Body: struct {
					Code      *int    `json:"code,omitempty"`
					Message   string  `json:"message"`
					RequestId *string `json:"request_id,omitempty"`
				}{
					Message: "Internal server error",
				},
			},
		}, nil
	}

	genRes := converters.FlatToGenArr(flats)

	return generated.GetHouseId200JSONResponse{
		Flats: genRes,
	}, nil
}
