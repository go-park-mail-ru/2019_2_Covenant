package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/playlist"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	. "2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	. "2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PlaylistHandler struct {
	BaseHandler
	PUsecase playlist.Usecase
}

func NewPlaylistHandler(pUC playlist.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *PlaylistHandler {
	return &PlaylistHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		PUsecase: pUC,
	}
}

func (ph *PlaylistHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/playlists", ph.CreatePlaylist(), ph.MManager.CheckAuth)
	e.GET("/api/v1/playlists", ph.GetPlaylists(), ph.MManager.CheckAuth)
	//e.DELETE("/api/v1/playlists/:id", ph.DeletePlaylust(), ph.MManager.CheckAuth)
	//e.GET("/api/v1/playlists/:id", ph.GetPlaylist(), ph.MManager.CheckAuth)
}

func (ph *PlaylistHandler) CreatePlaylist() echo.HandlerFunc {
	type Request struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			ph.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := ph.ReqReader.Read(c, request, nil); err != nil {
			ph.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		newPlaylist := models.NewPlaylist(request.Name, request.Description, sess.UserID)

		if err := ph.PUsecase.Store(newPlaylist); err != nil {
			ph.Logger.Log(c, "info", "Error while storing playlist.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"playlist": newPlaylist,
			},
		})
	}
}

func (ph *PlaylistHandler) GetPlaylists() echo.HandlerFunc {
	type Request struct {
		Count  uint64 `json:"count" validate:"required"`
		Offset uint64 `json:"offset"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			ph.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := ph.ReqReader.Read(c, request, nil); err != nil {
			ph.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		playlists, total, err := ph.PUsecase.Fetch(sess.UserID, request.Count, request.Offset)

		if err != nil {
			ph.Logger.Log(c, "error", "Error while fetching tracks.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"playlists": playlists,
				"total":  total,
			},
		})
	}
}
