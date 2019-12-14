package apiserver

import (
	_albumDelivery "2019_2_Covenant/internal/album/delivery"
	_albumUsecase "2019_2_Covenant/internal/album/usecase"
	"2019_2_Covenant/internal/app/storage"
	_artistDelivery "2019_2_Covenant/internal/artist/delivery"
	_artistUsecase "2019_2_Covenant/internal/artist/usecase"
	"2019_2_Covenant/internal/middlewares"
	_playlistDelivery "2019_2_Covenant/internal/playlist/delivery"
	_playlistUsecase "2019_2_Covenant/internal/playlist/usecase"
	_searchDelivery "2019_2_Covenant/internal/search/delivery"
	_searchUsecase "2019_2_Covenant/internal/search/usecase"
	_sessionDelivery "2019_2_Covenant/internal/session/delivery"
	_sessionUsecase "2019_2_Covenant/internal/session/usecase"
	_subscriptionDelivery "2019_2_Covenant/internal/subscriptions/delivery"
	_subscriptionUsecase "2019_2_Covenant/internal/subscriptions/usecase"
	_trackDelivery "2019_2_Covenant/internal/track/delivery"
	_trackUsecase "2019_2_Covenant/internal/track/usecase"
	_userDelivery "2019_2_Covenant/internal/user/delivery"
	_userUsecase "2019_2_Covenant/internal/user/usecase"
	"2019_2_Covenant/pkg/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

type APIServer struct {
	conf    *Config
	router  *echo.Echo
	storage storage.Storage
	logger  *logger.LogrusLogger
}

func NewAPIServer(conf *Config, st storage.Storage) *APIServer {
	return &APIServer{
		conf:    conf,
		router:  echo.New(),
		storage: st,
		logger:  logger.NewLogrusLogger(),
	}
}

func (api *APIServer) Start() error {
	if err := api.configureLogger(); err != nil {
		return nil
	}

	api.logger.L.Info("starting server...")

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
	playlistUsecase := _playlistUsecase.NewPlaylistUsecase(api.storage.Playlist())
	searchUsecase := _searchUsecase.NewSearchUsecase(api.storage.Track(), api.storage.Album(), api.storage.Artist())
	artistUsecase := _artistUsecase.NewArtistUsecase(api.storage.Artist())
	albumUsecase := _albumUsecase.NewAlbumUsecase(api.storage.Album())
	subscriptionUsecase := _subscriptionUsecase.NewSubscriptionUsecase(api.storage.Subscription())

	middlewareManager := middlewares.NewMiddlewareManager(userUsecase, sessionUsecase, api.logger)
	api.router.Use(middlewareManager.AccessLogMiddleware)
	api.router.Use(middlewareManager.PanicRecovering)
	api.router.Use(middlewareManager.CORSMiddleware)

	userHandler := _userDelivery.NewUserHandler(userUsecase, sessionUsecase, middlewareManager, api.logger)
	userHandler.Configure(api.router)

	trackHandler := _trackDelivery.NewTrackHandler(trackUsecase, middlewareManager, api.logger)
	trackHandler.Configure(api.router)

	sessionHandler := _sessionDelivery.NewSessionHandler(sessionUsecase, userUsecase, middlewareManager, api.logger)
	sessionHandler.Configure(api.router)

	playlistHandler := _playlistDelivery.NewPlaylistHandler(playlistUsecase, middlewareManager, api.logger)
	playlistHandler.Configure(api.router)

	searchHandler := _searchDelivery.NewSearchHandler(searchUsecase, userUsecase, middlewareManager, api.logger)
	searchHandler.Configure(api.router)

	artistHandler := _artistDelivery.NewArtistHandler(artistUsecase, middlewareManager, api.logger)
	artistHandler.Configure(api.router)

	subscriptionHandler := _subscriptionDelivery.NewSubscriptionHandler(subscriptionUsecase, middlewareManager, api.logger)
	subscriptionHandler.Configure(api.router)

	albumHandler := _albumDelivery.NewAlbumHandler(albumUsecase, middlewareManager, api.logger)
	albumHandler.Configure(api.router)
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

	api.logger.L.SetLevel(level)

	return nil
}

func (api *APIServer) Stop() {
	api.storage.Close()
}
