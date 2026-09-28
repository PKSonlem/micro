package flat

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/PKSonlem/micro/internal/deps"
	"github.com/PKSonlem/micro/internal/entity"
	"github.com/PKSonlem/micro/internal/generated"
	"github.com/PKSonlem/micro/internal/service/converters"
	flatusecases "github.com/PKSonlem/micro/internal/usecase/flat"
)

type Handler struct {
	usecase *flatusecases.Usecase
	logger  deps.Logger
}

func NewHandler(
	usecase *flatusecases.Usecase,
	logger deps.Logger,
) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) CreateFlat(ctx context.Context, request generated.PostFlatCreateRequestObject) (generated.PostFlatCreateResponseObject, error) {
	if request.Body == nil {
		return generated.PostFlatCreate400Response{}, nil
	}

	flat := converters.FlatFromGenCreate(generated.PostFlatCreateJSONBody(*request.Body))

	if err := validator.New().Struct(flat); err != nil {
		h.logger.Error(ctx, err)
		return generated.PostFlatCreate400Response{}, nil
	}

	res, err := h.usecase.HandleCreateFlat(ctx, flat)
	if err != nil {
		h.logger.Error(ctx, err)
		if errors.Is(err, entity.ErrorCreatingFlat) {
			return generated.PostFlatCreate400Response{}, nil
		}

		return generated.PostFlatCreate500JSONResponse{
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
	h.logger.Info(ctx, fmt.Sprintf("flat created with ID: %d", res.ID))

	return generated.PostFlatCreate200JSONResponse(genResp), nil
}
