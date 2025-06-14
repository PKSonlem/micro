package cmd

import (
	househandler "github.com/timurzdev/mentorship-test-task/internal/handler/house"
	"github.com/timurzdev/mentorship-test-task/internal/handler/middlewares/prometheus"
	"github.com/timurzdev/mentorship-test-task/internal/handler/server"
	"github.com/timurzdev/mentorship-test-task/internal/repository"
	houseusecases "github.com/timurzdev/mentorship-test-task/internal/usecase/house"
)

// контейнер внутренних зависимостей
type Internal struct {
	//external
	*Container

	repository     *repository.Repository
	testRepository *repository.Repository

	server *server.Server

	//handlers
	createHouseHandler *househandler.Handler

	//usecases
	createHouseUsecase *houseusecases.Usecase

	//middlewares
	prometheusMiddleware *prometheus.Middleware
}

func NewInternal(container *Container) *Internal {
	return &Internal{Container: container}
}

func (i *Internal) GetRepository() *repository.Repository {
	if i.repository == nil {
		i.repository = repository.NewRepository(i.GetPostgres())
	}

	return i.repository
}

func (i *Internal) GetServer() *server.Server {
	if i.server == nil {
		i.server = server.NewServer(
			i.GetLogger(),
			i.configuration.GetServerConfiguration().GetAddress(),
			i.GetHouseHandler(),
			i.GetPrometheusMiddleware(),
		)
	}

	return i.server
}

func (i *Internal) GetHouseHandler() *househandler.Handler {
	if i.createHouseHandler == nil {
		i.createHouseHandler = househandler.NewHandler(
			i.GetHouseUsecases(),
			i.GetLogger(),
		)
	}

	return i.createHouseHandler
}

func (i *Internal) GetHouseUsecases() *houseusecases.Usecase {
	if i.createHouseUsecase == nil {
		i.createHouseUsecase = houseusecases.NewUsecase(i.GetRepository())
	}

	return i.createHouseUsecase
}

func (i *Internal) GetPrometheusMiddleware() *prometheus.Middleware {
	if i.prometheusMiddleware == nil {
		i.prometheusMiddleware = prometheus.New(i.GetMetrics())
	}

	return i.prometheusMiddleware
}
