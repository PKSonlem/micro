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
	roles   deps.RolesReader
	usecase create_house.Usecase
}

func NewHandler(usecase create_house.Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	role, err := h.roles.GetRole(r.Context())
	if err != nil {
		handler.ErrorResponse(w, err)
		return
	}

	if !role.IsAdmin() {
		handler.ErrorResponse(w, handler.ErrNotFound)
		return
	}

	var bodyBytes []byte
	_, err = r.Body.Read(bodyBytes)
	defer r.Body.Close()

	genReq := generated.PostHouseCreateJSONBody{}
	err = json.Unmarshal(bodyBytes, &genReq)
	if err != nil {
		handler.ErrorResponse(w, err)
		return
	}

	house := converters.HouseFromGen(genReq)

	err = h.usecase.Handle(house)
	if err != nil {
		handler.ErrorResponse(w, err)
		return
	}

	handler.SuccessResponse(w, nil)
	return
}
