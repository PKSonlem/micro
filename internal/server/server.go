package server

import (
	"net/http"

	"github.com/timurzdev/mentorship-test-task/internal/generated"
	"github.com/timurzdev/mentorship-test-task/internal/handler/create_house"
)

type Server struct {
	createHouseHandler *create_house.Handler
}

func NewServer(chh *create_house.Handler) *Server {
	return &Server{createHouseHandler: chh}
}

// (GET /dummyLogin)
func (s *Server) GetDummyLogin(w http.ResponseWriter, r *http.Request, params generated.GetDummyLoginParams) {
}

// (POST /flat/create)
func (s *Server) PostFlatCreate(w http.ResponseWriter, r *http.Request) {}

// (POST /flat/update)
func (s *Server) PostFlatUpdate(w http.ResponseWriter, r *http.Request) {}

// (POST /house/create)
func (s *Server) PostHouseCreate(w http.ResponseWriter, r *http.Request) {
	s.createHouseHandler.Handle(w, r)
}

// (GET /house/{id})
func (s *Server) GetHouseId(w http.ResponseWriter, r *http.Request, id generated.HouseId) {}

// (POST /house/{id}/subscribe)
func (s *Server) PostHouseIdSubscribe(w http.ResponseWriter, r *http.Request, id generated.HouseId) {}

// (POST /login)
func (s *Server) PostLogin(w http.ResponseWriter, r *http.Request) {}

// (POST /register)
func (s *Server) PostRegister(w http.ResponseWriter, r *http.Request) {}
