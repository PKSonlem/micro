package house

import (
	"context"

	"github.com/PKSonlem/micro/internal/generated"
)

func (h *Handler) Subscribe(ctx context.Context, request generated.PostHouseIdSubscribeRequestObject) (generated.PostHouseIdSubscribeResponseObject, error) {
	if request.Body == nil {
		return generated.N400Response{}, nil
	}

	houseID := request.Id
	email := string(request.Body.Email)

	err := h.usecase.Subscribe(ctx, houseID, email)
	if err != nil {
		h.logger.Error(ctx, err)
		return generated.PostHouseIdSubscribe500JSONResponse{
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

	return generated.PostHouseIdSubscribe200Response{}, nil
}
