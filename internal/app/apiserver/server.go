package apiserver

import (
	"2019_2_Covenant/internal/app/storage"
	"2019_2_Covenant/internal/middlewares"
	_sessionDelivery "2019_2_Covenant/internal/session/delivery"
	_sessionUsecase "2019_2_Covenant/internal/session/usecase"
	_trackDelivery "2019_2_Covenant/internal/track/delivery"
	_trackUsecase "2019_2_Covenant/internal/track/usecase"
	_userDelivery "2019_2_Covenant/internal/user/delivery"
	_userUsecase "2019_2_Covenant/internal/user/usecase"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

type APIServer struct {
	conf   *Config
	router *echo.Echo
	storage storage.Storage
	logger *logrus.Logger
}

func NewAPIServer(conf *Config, st storage.Storage) *APIServer {
	return &APIServer{
		conf:   conf,
		router: echo.New(),
		storage: st,
		logger: logrus.New(),
	}
}

func (api *APIServer) Start() error {
	if err := api.configureLogger(); err != nil {
		return nil
	}

	api.logger.Info("starting server...")

	if err := api.configureStorage(); err != nil {
		return err
	}

	api.configureRouter()

	return api.router.Start(fmt.Sprintf("%s:%s", api.conf.Address, api.conf.Port))
}

func (api *APIServer) configureRouter() {
	api.router.GET("/docs/*", echoSwagger.WrapHandler)
	fs := http.FileServer(http.Dir("resources/"))
	api.router.GET("/resources/*", echo.WrapHandler(http.StripPrefix("/resources/", fs)))

	userUsecase := _userUsecase.NewUserUsecase(api.storage.User())
	sessionUsecase := _sessionUsecase.NewSessionUsecase(api.storage.Session())
	trackUsecase := _trackUsecase.NewTrackUsecase(api.storage.Track())

	middlewareManager := middlewares.NewMiddlewareManager(userUsecase, sessionUsecase, api.logger)
	api.router.Use(middlewareManager.AccessLogMiddleware)
	api.router.Use(middlewareManager.PanicRecovering)
	api.router.Use(middlewareManager.CORSMiddleware)

	userHandler := _userDelivery.NewUserHandler(userUsecase, sessionUsecase, middlewareManager, api.logger)
	userHandler.Configure(api.router)

	trackHandler := _trackDelivery.NewTrackHandler(trackUsecase, middlewareManager, api.logger)
	trackHandler.Configure(api.router)

	sessionHandler := _sessionDelivery.NewSessionHandler(sessionUsecase, middlewareManager, api.logger)
	sessionHandler.Configure(api.router)
}

func (api *APIServer) configureStorage() error {
	if err := api.storage.Open(); err != nil {
		return err
	}

	return nil
}

func (api *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(api.conf.LogLevel)

	if err != nil {
		return err
	}

	api.logger.SetLevel(level)

	return nil
}
