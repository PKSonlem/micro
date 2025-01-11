package create_house

import (
	"encoding/json"
	"net/http"

	"github.com/timurzdev/mentorship-test-task/internal/converters"
	"github.com/timurzdev/mentorship-test-task/internal/deps"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
	"github.com/timurzdev/mentorship-test-task/internal/handler"
	"github.com/timurzdev/mentorship-test-task/internal/usecase/create_house"
)

type Handler struct {
	// не имплементировано
	// roles   deps.RolesReader

	usecase *create_house.Usecase
	logger  deps.Logger
}

func NewHandler(
	usecase *create_house.Usecase,
	logger deps.Logger,
) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: добавить проверку роли
	// role, err := h.roles.GetRole(ctx)
	// if err != nil {
	// 	handler.ErrorResponse(w, err)
	// 	return
	// }
	//
	// if !role.IsAdmin() {
	// 	handler.ErrorResponse(w, handler.ErrNotFound)
	// 	return
	// }

	var genReq generated.PostHouseCreateJSONBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&genReq)
	if err != nil {
		h.logger.Error(ctx, err)
		handler.ErrorResponse(w, err)
		return
	}

	house := converters.HouseFromGen(genReq)

	res, err := h.usecase.Handle(ctx, house)
	if err != nil {
		h.logger.Error(ctx, err)
		handler.ErrorResponse(w, err)
		return
	}

	genResp := converters.HouseToGen(*res)
	bytes, _ := json.Marshal(genResp)

	handler.SuccessResponse(w, bytes)
}
