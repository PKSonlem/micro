package server

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/timurzdev/mentorship-test-task/internal/deps"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
	authhandler "github.com/timurzdev/mentorship-test-task/internal/handler/auth"
	flathandler "github.com/timurzdev/mentorship-test-task/internal/handler/flat"
	househandler "github.com/timurzdev/mentorship-test-task/internal/handler/house"
	"github.com/timurzdev/mentorship-test-task/internal/handler/middlewares/auth"
	"github.com/timurzdev/mentorship-test-task/internal/handler/middlewares/prometheus"
)

// Server реализует StrictServerInterface
type Server struct {
	logger         deps.Logger
	addr           string
	server         *http.Server
	houseHandler   *househandler.Handler
	flatHandler    *flathandler.Handler
	authHandler    *authhandler.Handler
	promMiddleware *prometheus.Middleware
	authMiddleware *auth.Middleware
}

func NewServer(
	logger deps.Logger,
	addr string,
	houseHandler *househandler.Handler,
	flatHandler *flathandler.Handler,
	authHandler *authhandler.Handler,
	promMiddleware *prometheus.Middleware,
	authMiddleware *auth.Middleware,
) *Server {
	s := &Server{
		logger:         logger,
		addr:           addr,
		houseHandler:   houseHandler,
		flatHandler:    flatHandler,
		authHandler:    authHandler,
		promMiddleware: promMiddleware,
		authMiddleware: authMiddleware,
	}

	// Создаем мультиплексор
	mux := http.NewServeMux()

	// Регистрируем обработчики для метрик
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Создаем strict handler для API
	strictHandler := generated.NewStrictHandler(s, nil)

	// Применяем middleware в правильном порядке:
	// 1. Prometheus middleware для всех запросов
	// 2. Auth middleware (которая уже содержит логику исключений)
	// 3. Generated handler
	handler := generated.Handler(strictHandler)
	handler = s.authMiddleware.Apply(handler)
	handler = s.promMiddleware.Handle(handler)

	// Регистрируем обработчик API
	mux.Handle("/", handler)

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return s
}

func (s *Server) Run(ctx context.Context) {
	s.logger.Info(ctx, "server starting", "address", s.addr)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error(ctx, err)
	}
}

// Shutdown останавливает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Реализация StrictServerInterface

// PostHouseCreate делегирует обработку в house handler
func (s *Server) PostHouseCreate(ctx context.Context, request generated.PostHouseCreateRequestObject) (generated.PostHouseCreateResponseObject, error) {
	return s.houseHandler.CreateHouse(ctx, request)
}

// GetDummyLogin делегирует обработку в auth handler
func (s *Server) GetDummyLogin(ctx context.Context, request generated.GetDummyLoginRequestObject) (generated.GetDummyLoginResponseObject, error) {
	return s.authHandler.GetDummyLogin(ctx, request)
}

// PostFlatCreate
func (s *Server) PostFlatCreate(ctx context.Context, request generated.PostFlatCreateRequestObject) (generated.PostFlatCreateResponseObject, error) {
	return s.flatHandler.CreateFlat(ctx, request)
}

// PostFlatUpdate
func (s *Server) PostFlatUpdate(ctx context.Context, request generated.PostFlatUpdateRequestObject) (generated.PostFlatUpdateResponseObject, error) {
	return s.flatHandler.UpdateModeratorFlat(ctx, request)
}

// GetHouseId
func (s *Server) GetHouseId(ctx context.Context, request generated.GetHouseIdRequestObject) (generated.GetHouseIdResponseObject, error) {
	return s.houseHandler.GetHouseFlats(ctx, request)
}

// PostHouseIdSubscribe - заглушка
func (s *Server) PostHouseIdSubscribe(ctx context.Context, request generated.PostHouseIdSubscribeRequestObject) (generated.PostHouseIdSubscribeResponseObject, error) {
	return generated.PostHouseIdSubscribe500JSONResponse{
		N5xxJSONResponse: generated.N5xxJSONResponse{
			Body: struct {
				Code      *int    `json:"code,omitempty"`
				Message   string  `json:"message"`
				RequestId *string `json:"request_id,omitempty"`
			}{
				Message: "Not implemented",
			},
		},
	}, nil
}

// PostLogin - заглушка
func (s *Server) PostLogin(ctx context.Context, request generated.PostLoginRequestObject) (generated.PostLoginResponseObject, error) {
	return generated.PostLogin500JSONResponse{
		N5xxJSONResponse: generated.N5xxJSONResponse{
			Body: struct {
				Code      *int    `json:"code,omitempty"`
				Message   string  `json:"message"`
				RequestId *string `json:"request_id,omitempty"`
			}{
				Message: "Not implemented",
			},
		},
	}, nil
}

// PostRegister - заглушка
func (s *Server) PostRegister(ctx context.Context, request generated.PostRegisterRequestObject) (generated.PostRegisterResponseObject, error) {
	return generated.PostRegister500JSONResponse{
		N5xxJSONResponse: generated.N5xxJSONResponse{
			Body: struct {
				Code      *int    `json:"code,omitempty"`
				Message   string  `json:"message"`
				RequestId *string `json:"request_id,omitempty"`
			}{
				Message: "Not implemented",
			},
		},
	}, nil
}
