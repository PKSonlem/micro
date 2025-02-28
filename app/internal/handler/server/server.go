package server

import (
	"context"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/timurzdev/mentorship-test-task/internal/deps"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
	"github.com/timurzdev/mentorship-test-task/internal/handler/create_house"
	"github.com/timurzdev/mentorship-test-task/internal/handler/middlewares/prometheus"
)

type Server struct {
	logger             deps.Logger
	address            string
	createHouseHandler *create_house.Handler

	prometheusMiddleware *prometheus.Middleware
}

func NewServer(
	log deps.Logger,
	address string,
	chh *create_house.Handler,
	prometheusMiddleware *prometheus.Middleware,
) *Server {
	return &Server{
		address:              address,
		createHouseHandler:   chh,
		logger:               log,
		prometheusMiddleware: prometheusMiddleware,
	}
}

func (s *Server) Run(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("GET /metrics", promhttp.Handler())

	// биндим нашу структуру Server к роутам
	h := generated.HandlerWithOptions(s, generated.StdHTTPServerOptions{
		BaseRouter: mux,
		Middlewares: []generated.MiddlewareFunc{
			s.prometheusMiddleware.Handle,
		},
	})

	srv := &http.Server{
		Handler: h,
		Addr:    s.address,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	// старт http сервера
	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Error(ctx, err)
	}
}

// (POST /house/create)
func (s *Server) PostHouseCreate(w http.ResponseWriter, r *http.Request) {
	s.createHouseHandler.Handle(w, r)
}

// (GET /dummyLogin)
func (s *Server) GetDummyLogin(w http.ResponseWriter, r *http.Request, params generated.GetDummyLoginParams) {
	//not implemented
}

// (POST /flat/create)
func (s *Server) PostFlatCreate(w http.ResponseWriter, r *http.Request) {
	//not implemented
}

// (POST /flat/update)
func (s *Server) PostFlatUpdate(w http.ResponseWriter, r *http.Request) {
	//not implemented
}

// (GET /house/{id})
func (s *Server) GetHouseId(w http.ResponseWriter, r *http.Request, id generated.HouseId) {
	//not implemented
}

// (POST /house/{id}/subscribe)
func (s *Server) PostHouseIdSubscribe(w http.ResponseWriter, r *http.Request, id generated.HouseId) {
	//not implemented
}

// (POST /login)
func (s *Server) PostLogin(w http.ResponseWriter, r *http.Request) {
	//not implemented
}

// (POST /register)
func (s *Server) PostRegister(w http.ResponseWriter, r *http.Request) {
	//not implemented
}
