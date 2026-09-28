package house

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/PKSonlem/micro/internal/deps"
	"github.com/PKSonlem/micro/internal/entity"
	"github.com/PKSonlem/micro/internal/generated"
	"github.com/PKSonlem/micro/internal/service/converters"
	"github.com/PKSonlem/micro/internal/service/roles"
	houseusecases "github.com/PKSonlem/micro/internal/usecase/house"
)

type Handler struct {
	roleProvider *roles.RolesProvider
	usecase      *houseusecases.Usecase
	logger       deps.Logger
}

func NewHandler(
	roleProvider *roles.RolesProvider,
	usecase *houseusecases.Usecase,
	logger deps.Logger,
) *Handler {
	return &Handler{
		roleProvider: roleProvider,
		usecase:      usecase,
		logger:       logger,
	}
}

// CreateHouse обрабатывает создание дома с типизированными объектами
func (h *Handler) CreateHouse(ctx context.Context, request generated.PostHouseCreateRequestObject) (generated.PostHouseCreateResponseObject, error) {
	if request.Body == nil {
		return generated.PostHouseCreate400Response{}, nil
	}

	house := converters.HouseFromGen(generated.PostHouseCreateJSONBody(*request.Body))

	// Валидация
	if err := validator.New().Struct(house); err != nil {
		h.logger.Error(ctx, err)
		return generated.PostHouseCreate400Response{}, nil
	}

	res, err := h.usecase.Handle(ctx, house)
	if err != nil {
		h.logger.Error(ctx, err)

		// Все ошибки создания дома считаем ошибками валидации (400)
		// В реальном приложении мы бы проверяли уникальность адреса до попытки вставки
		if errors.Is(err, entity.ErrorCreatingHouse) {
			return generated.PostHouseCreate400Response{}, nil
		}

		// Для неожиданных ошибок возвращаем 500
		return generated.PostHouseCreate500JSONResponse{
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

	genResp := converters.HouseToGen(*res)
	h.logger.Info(ctx, fmt.Sprintf("house created with ID: %d", res.ID))

	return generated.PostHouseCreate200JSONResponse(genResp), nil
}
