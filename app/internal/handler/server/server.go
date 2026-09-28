package server

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/PKSonlem/micro/internal/deps"
	"github.com/PKSonlem/micro/internal/generated"

	authhandler "github.com/PKSonlem/micro/internal/handler/auth"
	flathandler "github.com/PKSonlem/micro/internal/handler/flat"
	househandler "github.com/PKSonlem/micro/internal/handler/house"
	"github.com/PKSonlem/micro/internal/handler/middlewares/auth"
	"github.com/PKSonlem/micro/internal/handler/middlewares/prometheus"
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

// PostHouseCreate
func (s *Server) PostHouseCreate(ctx context.Context, request generated.PostHouseCreateRequestObject) (generated.PostHouseCreateResponseObject, error) {
	return s.houseHandler.CreateHouse(ctx, request)
}

// GetDummyLogin
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

// PostHouseIdSubscribe
func (s *Server) PostHouseIdSubscribe(ctx context.Context, request generated.PostHouseIdSubscribeRequestObject) (generated.PostHouseIdSubscribeResponseObject, error) {
	return s.houseHandler.Subscribe(ctx, request)
}

// PostLogin
func (s *Server) PostLogin(ctx context.Context, request generated.PostLoginRequestObject) (generated.PostLoginResponseObject, error) {
	return s.authHandler.PostLogin(ctx, request)
}

// PostRegister
func (s *Server) PostRegister(ctx context.Context, request generated.PostRegisterRequestObject) (generated.PostRegisterResponseObject, error) {
	return s.authHandler.PostRegister(ctx, request)
}
