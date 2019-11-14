package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/track"
	"2019_2_Covenant/internal/vars"
	"2019_2_Covenant/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type TrackHandler struct {
	TUsecase track.Usecase
	MManager middlewares.MiddlewareManager
	Logger   *logger.LogrusLogger
}

func NewTrackHandler(tUC track.Usecase,
	mManager middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *TrackHandler {
	return &TrackHandler{
		TUsecase: tUC,
		MManager: mManager,
		Logger:   logger,
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
// @Failure 400 object vars.ResponseError
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
// @Router /api/v1/tracks/popular [post]
func (th *TrackHandler) GetPopularTracks() echo.HandlerFunc {
	return func(c echo.Context) error {
		tracks, err := th.TUsecase.FetchPopular(25)

		if err != nil {
			th.Logger.Log(c, "error", "Error while fetching tracks.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		for _, item := range tracks {
			start := strings.Index(item.Duration, "T")
			end := strings.Index(item.Duration, "Z")
			item.Duration = item.Duration[start+1 : end]
		}

		return c.JSON(http.StatusOK, vars.Response{Body: tracks})
	}
}

func (th *TrackHandler) AddToFavourites() echo.HandlerFunc {
	type DataToAdd struct {
		TrackID uint64 `json:"track_id"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			th.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		data := &DataToAdd{}

		if err := c.Bind(&data); err != nil {
			th.Logger.Log(c, "error","Can't read request body.")
			return c.JSON(http.StatusUnprocessableEntity, vars.ResponseError{Error: err.Error()})
		}

		if err := th.TUsecase.StoreFavourite(sess.UserID, data.TrackID); err != nil {
			th.Logger.Log(c, "error", "Error while storing favourite track.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "success",
		})
	}
}

func (th *TrackHandler) RemoveFavourite() echo.HandlerFunc {
	type DataToRemove struct {
		TrackID uint64 `json:"track_id"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			th.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		data := &DataToRemove{}

		if err := c.Bind(&data); err != nil {
			th.Logger.Log(c, "error","Can't read request body.")
			return c.JSON(http.StatusUnprocessableEntity, vars.ResponseError{Error: err.Error()})
		}

		if err := th.TUsecase.RemoveFavourite(sess.UserID, data.TrackID); err != nil {
			th.Logger.Log(c, "error", "Error while remove favourite track.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "success",
		})
	}
}

func (th *TrackHandler) GetFavourites() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			th.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		tracks, err := th.TUsecase.FetchFavourites(sess.UserID, 25)

		if err != nil {
			th.Logger.Log(c, "error", "Error while fetching tracks.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		for _, item := range tracks {
			start := strings.Index(item.Duration, "T")
			end := strings.Index(item.Duration, "Z")
			item.Duration = item.Duration[start+1 : end]
		}

		return c.JSON(http.StatusOK, vars.Response{Body: tracks})
	}
}
