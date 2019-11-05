package delivery

import (
	"2019_2_Covenant/internal/middleware"
	"2019_2_Covenant/internal/track"
	"2019_2_Covenant/internal/vars"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type TrackHandler struct {
	TUsecase track.Usecase
	MManager middleware.MiddlewareManager
}

func NewTrackHandler(tUC track.Usecase, mManager middleware.MiddlewareManager) *TrackHandler {
	return &TrackHandler{
		TUsecase: tUC,
		MManager: mManager,
	}
}

func (th *TrackHandler) Configure(e *echo.Echo) {
	e.GET("/api/v1/tracks/popular", th.GetPopularTracks())
}

type ResponseError struct {
	Error string `json:"error"`
}

type Response struct {
	Body interface{} `json:"body"`
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
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		for _, track := range tracks {
			start := strings.Index(track.Duration, "T")
			end := strings.Index(track.Duration, "Z")
			track.Duration = track.Duration[start+1 : end]
		}

		return c.JSON(http.StatusOK, Response{tracks})
	}
}