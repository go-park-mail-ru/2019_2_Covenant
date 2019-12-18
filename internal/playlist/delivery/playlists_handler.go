package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/playlist"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	. "2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	"2019_2_Covenant/tools/time_parser"
	. "2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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
	e.POST("/api/v1/playlists", ph.CreatePlaylist(), ph.MManager.CheckAuthStrictly)
	e.GET("/api/v1/playlists", ph.GetPlaylists(), ph.MManager.CheckAuthStrictly)
	e.GET("/api/v1/playlists/:id", ph.GetSinglePlaylist(), ph.MManager.CheckAuthStrictly)
	e.GET("/api/v1/playlists/:id/tracks", ph.GetTracksFromPlaylist(), ph.MManager.CheckAuthStrictly)
	e.DELETE("/api/v1/playlists/:id", ph.DeletePlaylist(), ph.MManager.CheckAuthStrictly)
	e.POST("/api/v1/playlists/:id/tracks", ph.AddToPlaylist(), ph.MManager.CheckAuthStrictly)
	e.DELETE("/api/v1/playlists/:playlist_id/tracks/:track_id", ph.RemoveFromPlaylist(), ph.MManager.CheckAuthStrictly)
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
		Count  uint64 `query:"count" validate:"required"`
		Offset uint64 `query:"offset"`
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
			ph.Logger.Log(c, "error", "Error while fetching playlists.", err)
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

func (ph *PlaylistHandler) DeletePlaylist() echo.HandlerFunc {
	return func(c echo.Context) error {
		pID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ph.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		if err := ph.PUsecase.DeleteByID(uint64(pID)); err != nil {
			ph.Logger.Log(c, "info", "Error while remove playlist.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ph *PlaylistHandler) AddToPlaylist() echo.HandlerFunc {
	type Request struct {
		TrackID uint64 `json:"track_id" validate:"required"`
	}

	return func(c echo.Context) error {
		pID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ph.Logger.Log(c, "error", "Atoi error.", err.Error())
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

		if err := ph.PUsecase.AddToPlaylist(uint64(pID), request.TrackID); err != nil {
			ph.Logger.Log(c, "error", "Error while adding track to playlist.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrAlreadyExist.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ph *PlaylistHandler) RemoveFromPlaylist() echo.HandlerFunc {
	return func(c echo.Context) error {
		pID, err1 := strconv.Atoi(c.Param("playlist_id"))
		tID, err2 := strconv.Atoi(c.Param("track_id"))

		if err1 != nil || err2 != nil {
			ph.Logger.Log(c, "error", "Atoi error.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		if err := ph.PUsecase.RemoveFromPlaylist(uint64(pID), uint64(tID)); err != nil {
			ph.Logger.Log(c, "info", "Error while remove playlist.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ph *PlaylistHandler) GetSinglePlaylist() echo.HandlerFunc {
	return func(c echo.Context) error {
		pID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ph.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		p, amountOfTracks, err := ph.PUsecase.GetSinglePlaylist(uint64(pID))

		if err != nil {
			ph.Logger.Log(c, "info", "Error while getting playlist.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"playlist": p,
				"amount_of_tracks": amountOfTracks,
			},
		})
	}
}

func (ph *PlaylistHandler) GetTracksFromPlaylist() echo.HandlerFunc {
	return func(c echo.Context) error {
		pID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ph.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		// TODO: Strictly

		var authID uint64
		if sess, ok := c.Get("session").(*models.Session); ok {
			authID = sess.UserID
		}

		tracks, err := ph.PUsecase.GetTracksFrom(uint64(pID), authID)

		if err != nil {
			ph.Logger.Log(c, "error", "Error while fetching tracks.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		for _, item := range tracks { item.Duration = time_parser.GetDuration(item.Duration) }

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"tracks": tracks,
			},
		})
	}
}
