package apiserver

import (
	"2019_2_Covenant/internal/app/storage"
	_sessionUsecase "2019_2_Covenant/internal/session/usecase"
	"2019_2_Covenant/internal/user/delivery"
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
	api.configureRouter()

	if err := api.configureStorage(); err != nil {
		return err
	}

	return api.router.Start(fmt.Sprintf("%s:%s", api.conf.Address, api.conf.Port))
}

func (api *APIServer) configureRouter() {
	api.router.GET("/docs/*", echoSwagger.WrapHandler)

	userUsecase := _userUsecase.NewUserUsecase(api.storage.User())
	sessionUsecase := _sessionUsecase.NewSessionUsecase(api.storage.Session())

	handler := delivery.NewUserHandler(userUsecase, sessionUsecase)
	handler.Configure(api.router)
}

func (api *APIServer) configureStorage() error {
	if err := api.storage.Open(); err != nil {
		return err
	}

	return nil
}
