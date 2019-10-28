package apiserver

import (
	"2019_2_Covenant/internal/app/storage"
	_sessionRepo "2019_2_Covenant/internal/session/repository"
	_sessionUsecase "2019_2_Covenant/internal/session/usecase"
	"2019_2_Covenant/internal/user/delivery"
	_userRepo "2019_2_Covenant/internal/user/repository"
	_userUsecase "2019_2_Covenant/internal/user/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type APIServer struct {
	conf   *Config
	router *echo.Echo
	storage *storage.Storage
}

func NewAPIServer(conf *Config) *APIServer {
	return &APIServer{
		conf:   conf,
		router: echo.New(),
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

	userStorage := _userRepo.NewUserStorage()
	userUsecase := _userUsecase.NewUserUsecase(userStorage)
	sessionStorage := _sessionRepo.NewSessionStorage()
	sessionUsecase := _sessionUsecase.NewSessionUsecase(sessionStorage)

	handler := delivery.NewUserHandler(userUsecase, sessionUsecase)
	handler.Configure(api.router)
}

func (api *APIServer) configureStorage() error {
	st := storage.NewStorage(api.conf.Storage)

	if err := st.Open(); err != nil {
		return err
	}

	api.storage = st

	return nil
}
