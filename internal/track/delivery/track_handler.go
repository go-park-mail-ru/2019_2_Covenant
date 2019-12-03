package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/track"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	"2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	. "2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type TrackHandler struct {
	base_handler.BaseHandler
	TUsecase track.Usecase
}

func NewTrackHandler(tUC track.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *TrackHandler {
	return &TrackHandler{
		BaseHandler: base_handler.BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		TUsecase: tUC,
	}
}

func (th *TrackHandler) Configure(e *echo.Echo) {
	e.GET("/api/v1/tracks/popular", th.GetPopularTracks())

	e.GET("/api/v1/tracks/favourite", th.GetFavourites(), th.MManager.CheckAuth)
	e.POST("/api/v1/tracks/favourite", th.AddToFavourites(), th.MManager.CheckAuth)
	e.DELETE("/api/v1/tracks/favourite", th.RemoveFavourite(), th.MManager.CheckAuth)
}

// @Tags Track
// @Summary Get Popular Tracks Route
// @Description Getting popular tracks
// @ID get-popular-tracks
// @Accept json
// @Produce json
// @Success 200 object models.Track
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/tracks/popular [post]
func (th *TrackHandler) GetPopularTracks() echo.HandlerFunc {
	type Request struct {
		Count  uint64 `query:"count" validate:"required"`
		Offset uint64 `query:"offset"`
	}

	return func(c echo.Context) error {
		request := &Request{}

		if err := th.ReqReader.Read(c, request, nil); err != nil {
			th.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		tracks, err := th.TUsecase.FetchPopular(request.Count, request.Offset)

		if err != nil {
			th.Logger.Log(c, "error", "Error while fetching tracks.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		for _, item := range tracks {
			start := strings.Index(item.Duration, "T")
			end := strings.Index(item.Duration, "Z")
			item.Duration = item.Duration[start+1 : end]
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"tracks": tracks,
			},
		})
	}
}

func (th *TrackHandler) AddToFavourites() echo.HandlerFunc {
	type Request struct {
		TrackID uint64 `json:"track_id" validate:"required"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			th.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := th.ReqReader.Read(c, request, nil); err != nil {
			th.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := th.TUsecase.StoreFavourite(sess.UserID, request.TrackID); err != nil {
			th.Logger.Log(c, "error", "Error while storing favourite track.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrAlreadyExist.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (th *TrackHandler) RemoveFavourite() echo.HandlerFunc {
	type Request struct {
		TrackID uint64 `json:"track_id" validate:"required"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			th.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := th.ReqReader.Read(c, request, nil); err != nil {
			th.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := th.TUsecase.RemoveFavourite(sess.UserID, request.TrackID); err != nil {
			th.Logger.Log(c, "error", "Error while remove favourite track.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (th *TrackHandler) GetFavourites() echo.HandlerFunc {
	type Request struct {
		Count  uint64 `query:"count" validate:"required"`
		Offset uint64 `query:"offset"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			th.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := th.ReqReader.Read(c, request, nil); err != nil {
			th.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		tracks, total, err := th.TUsecase.FetchFavourites(sess.UserID, request.Count, request.Offset)

		if err != nil {
			th.Logger.Log(c, "error", "Error while fetching tracks.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		for _, item := range tracks {
			start := strings.Index(item.Duration, "T")
			end := strings.Index(item.Duration, "Z")
			item.Duration = item.Duration[start+1 : end]
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"tracks": tracks,
				"total":  total,
			},
		})
	}
}
