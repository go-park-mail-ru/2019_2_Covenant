package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/track"
	"2019_2_Covenant/internal/vars"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type TrackHandler struct {
	TUsecase track.Usecase
	MManager middlewares.MiddlewareManager
}

func NewTrackHandler(tUC track.Usecase, mManager middlewares.MiddlewareManager) *TrackHandler {
	return &TrackHandler{
		TUsecase: tUC,
		MManager: mManager,
	}
}

func (th *TrackHandler) Configure(e *echo.Echo) {
	e.GET("/api/v1/tracks/popular", th.GetPopularTracks())
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
	return func(c echo.Context) error {
		tracks, err := th.TUsecase.Fetch(25)

		if err != nil {
			fmt.Println(err)
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
