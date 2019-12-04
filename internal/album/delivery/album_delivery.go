package delivery

import (
	"2019_2_Covenant/internal/album"
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	. "2019_2_Covenant/tools/base_handler"
	"github.com/labstack/echo/v4"
)

type AlbumHandler struct {
	BaseHandler
	AUsecase album.Usecase
}

func NewAlbumHandler(aUC album.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *AlbumHandler {
	return &AlbumHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		AUsecase: aUC,
	}
}

func (ah *AlbumHandler) Configure(e *echo.Echo) {
	//e.DELETE("/api/v1/albums/:id", ah.DeleteAlbum(), ah.MManager.CheckAuth, ah.MManager.CheckAdmin)
	//e.PUT("/api/v1/albums/:id", ah.UpdateAlbum(), ah.MManager.CheckAuth, ah.MManager.CheckAdmin)
	//e.GET("/api/v1/albums", ah.GetAlbums())
}
