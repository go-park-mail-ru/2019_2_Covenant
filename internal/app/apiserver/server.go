package apiserver

import (
	"2019_2_Covenant/internal/app/storage"
	"2019_2_Covenant/internal/middlewares"
	_sessionUsecase "2019_2_Covenant/internal/session/usecase"
	_trackDelivery "2019_2_Covenant/internal/track/delivery"
	_trackUsecase "2019_2_Covenant/internal/track/usecase"
	_userDelivery "2019_2_Covenant/internal/user/delivery"
	_userUsecase "2019_2_Covenant/internal/user/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type APIServer struct {
	conf   *Config
	router *echo.Echo
	storage storage.Storage
}

func NewAPIServer(conf *Config, st storage.Storage) *APIServer {
	return &APIServer{
		conf:   conf,
		router: echo.New(),
		storage: st,
	}
}

func (api *APIServer) Start() error {
	if err := api.configureStorage(); err != nil {
		return err
	}

	api.configureRouter()

	return api.router.Start(fmt.Sprintf("%s:%s", api.conf.Address, api.conf.Port))
}

func (api *APIServer) configureRouter() {
	api.router.GET("/docs/*", echoSwagger.WrapHandler)

	userUsecase := _userUsecase.NewUserUsecase(api.storage.User())
	sessionUsecase := _sessionUsecase.NewSessionUsecase(api.storage.Session())
	trackUsecase := _trackUsecase.NewTrackUsecase(api.storage.Track())

	middlewareManager := middlewares.NewMiddlewareManager(userUsecase, sessionUsecase)
	api.router.Use(middlewareManager.PanicRecovering)

	userHandler := _userDelivery.NewUserHandler(userUsecase, sessionUsecase, middlewareManager)
	userHandler.Configure(api.router)
	trackHandler := _trackDelivery.NewTrackHandler(trackUsecase, middlewareManager)
	trackHandler.Configure(api.router)
}

func (api *APIServer) configureStorage() error {
	if err := api.storage.Open(); err != nil {
		return err
	}

	return nil
}
