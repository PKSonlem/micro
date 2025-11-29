package flat

import (
	"context"
	"errors"
	"fmt"

	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
	"github.com/timurzdev/mentorship-test-task/internal/service/converters"
	"github.com/timurzdev/mentorship-test-task/internal/service/roles"
)

func (h *Handler) UpdateModeratorFlat(ctx context.Context, request generated.PostFlatUpdateRequestObject) (generated.PostFlatUpdateResponseObject, error) {
	if request.Body == nil {
		return generated.PostFlatUpdate400Response{}, nil
	}

	role, err := roles.NewProvider().GetRole(ctx)
	if err != nil {
		h.logger.Error(ctx, err)
		return generated.PostFlatUpdate400Response{}, nil
	}

	if !role.IsModerator() {
		h.logger.Error(ctx, fmt.Errorf("user is not a moderator"))
		return generated.PostFlatUpdate401Response{}, nil
	}

	flat := converters.FlatFromGenUpdate(generated.PostFlatUpdateJSONBody(*request.Body))

	flatID := flat.ID
	status := flat.Status

	if !entity.IsValidStatus(status) {
		h.logger.Error(ctx, fmt.Errorf("invalid status: %s", status))
		return generated.PostFlatUpdate400Response{}, nil
	}

	res, err := h.usecase.HandleUpdateStatus(ctx, flatID, status)
	if err != nil {
		h.logger.Error(ctx, err)
		if errors.Is(err, entity.ErrorUpdateModeratorFlat) {
			return generated.PostFlatUpdate400Response{}, nil
		}

		return generated.PostFlatUpdate500JSONResponse{
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

	genResp := converters.FlatToGen(*res)
	h.logger.Info(ctx, fmt.Sprintf("flat updated with ID: %d, new status: %s", res.ID, res.Status))

	return generated.PostFlatUpdate200JSONResponse(genResp), nil
}
