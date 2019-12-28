package api

import (
	"2019_2_Covenant/pkg/album"
	_albumDelivery "2019_2_Covenant/pkg/album/delivery"
	_albumUsecase "2019_2_Covenant/pkg/album/usecase"
	"2019_2_Covenant/pkg/artist"
	_artistDelivery "2019_2_Covenant/pkg/artist/delivery"
	_artistUsecase "2019_2_Covenant/pkg/artist/usecase"
	files "2019_2_Covenant/pkg/file_processor"
	"2019_2_Covenant/pkg/likes"
	_likesDelivery "2019_2_Covenant/pkg/likes/delivery"
	_likesUsecase "2019_2_Covenant/pkg/likes/usecase"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/middlewares"
	"2019_2_Covenant/pkg/playlist"
	_playlistDelivery "2019_2_Covenant/pkg/playlist/delivery"
	_playlistUsecase "2019_2_Covenant/pkg/playlist/usecase"
	_searchDelivery "2019_2_Covenant/pkg/search/delivery"
	_searchUsecase "2019_2_Covenant/pkg/search/usecase"
	_sessionDelivery "2019_2_Covenant/pkg/session/delivery"
	session "2019_2_Covenant/pkg/session/repository"
	_sessionUsecase "2019_2_Covenant/pkg/session/usecase"
	"2019_2_Covenant/pkg/subscriptions"
	_subscriptionDelivery "2019_2_Covenant/pkg/subscriptions/delivery"
	_subscriptionUsecase "2019_2_Covenant/pkg/subscriptions/usecase"
	"2019_2_Covenant/pkg/track"
	_trackDelivery "2019_2_Covenant/pkg/track/delivery"
	_trackUsecase "2019_2_Covenant/pkg/track/usecase"
	_userDelivery "2019_2_Covenant/pkg/user/delivery"
	user "2019_2_Covenant/pkg/user/repository"
	_userUsecase "2019_2_Covenant/pkg/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

type APIStorage interface {
	Track() track.Repository
	Playlist() playlist.Repository
	Album() album.Repository
	Artist() artist.Repository
	Subscription() subscriptions.Repository
	Like() likes.Repository
}

type AuthStorage interface {
	User() user.UsersClient
	Session() session.SessionsClient
}

type APIServer struct {
	router         *echo.Echo
	logger         *logger.LogrusLogger
	apiStorage     APIStorage
	authStorage    AuthStorage
	fileRepository files.Repository
}

func NewAPIServer(apiStorage APIStorage, authStorage AuthStorage, fileRepository files.Repository) *APIServer {
	return &APIServer{
		router:         echo.New(),
		apiStorage:     apiStorage,
		logger:         logger.NewLogrusLogger(),
		authStorage:    authStorage,
		fileRepository: fileRepository,
	}
}

func (api *APIServer) Start(address string, logLevel string, filesDir string) error {
	if err := api.configureLogger(logLevel); err != nil {
		return nil
	}

	api.logger.L.Info("starting server...")

	api.configureRouter(filesDir)

	return api.router.Start(address)
}

func (api *APIServer) configureRouter(filesDir string) {
	api.router.GET("/docs/*", echoSwagger.WrapHandler)

	fs := http.FileServer(http.Dir(filesDir))
	api.router.GET("/resources/*", echo.WrapHandler(http.StripPrefix("/resources/", fs)))

	userUsecase := _userUsecase.NewUserUsecase(api.authStorage.User(), api.fileRepository)
	sessionUsecase := _sessionUsecase.NewSessionUsecase(api.authStorage.Session())
	trackUsecase := _trackUsecase.NewTrackUsecase(api.apiStorage.Track(), api.fileRepository)
	playlistUsecase := _playlistUsecase.NewPlaylistUsecase(api.apiStorage.Playlist())
	searchUsecase := _searchUsecase.NewSearchUsecase(api.apiStorage.Track(), api.apiStorage.Album(), api.apiStorage.Artist())
	artistUsecase := _artistUsecase.NewArtistUsecase(api.apiStorage.Artist(), api.fileRepository)
	albumUsecase := _albumUsecase.NewAlbumUsecase(api.apiStorage.Album(), api.fileRepository)
	subscriptionUsecase := _subscriptionUsecase.NewSubscriptionUsecase(api.apiStorage.Subscription())
	likesUsecase := _likesUsecase.NewLikesUsecase(api.apiStorage.Like())

	middlewareManager := middlewares.NewMiddlewareManager(userUsecase, sessionUsecase, api.logger)
	api.router.Use(middlewareManager.AccessLogMiddleware)
	api.router.Use(middlewareManager.PanicRecovering)
	api.router.Use(middlewareManager.CORSMiddleware)

	userHandler := _userDelivery.NewUserHandler(userUsecase, sessionUsecase, playlistUsecase, middlewareManager, api.logger)
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

	subscriptionHandler := _subscriptionDelivery.NewSubscriptionHandler(subscriptionUsecase, userUsecase, middlewareManager, api.logger)
	subscriptionHandler.Configure(api.router)

	albumHandler := _albumDelivery.NewAlbumHandler(albumUsecase, middlewareManager, api.logger)
	albumHandler.Configure(api.router)

	likesHandler := _likesDelivery.NewLikesHandler(likesUsecase, userUsecase, middlewareManager, api.logger)
	likesHandler.Configure(api.router)
}

func (api *APIServer) configureLogger(logLevel string) error {
	level, err := logrus.ParseLevel(logLevel)

	if err != nil {
		return err
	}

	api.logger.L.SetLevel(level)

	return nil
}

func (api *APIServer) Stop() {
}
