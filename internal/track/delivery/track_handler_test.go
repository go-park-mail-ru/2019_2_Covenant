package delivery

import (
	. "2019_2_Covenant/internal/models"
	mock "2019_2_Covenant/internal/track/mocks"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

//go:generate mockgen -source=../usecase.go -destination=../mocks/mock_usecase.go -package=mock

func TestTrackHandler_GetPopularTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TUsecase := mock.NewMockRepository(ctrl)

	handler := TrackHandler{TUsecase: TUsecase}

	count := uint64(25)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/popular")

		tracks := []*Track{
			{ID: 1, Name: "Still loving you", Duration: "2019-10-31T00:06:28Z"},
		}
		TUsecase.EXPECT().Fetch(count).Return(tracks, nil)
		err := handler.GetPopularTracks()(c)

		if err != nil {
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"body":[{"name":"Still loving you","duration":"00:06:28","photo":"","artist":"","album":"","path":""}]}` {
			t1.Fail()
		}
	})

	t.Run("Error", func(t2 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/popular")

		TUsecase.EXPECT().Fetch(count).Return(nil, fmt.Errorf("some error"))
		err := handler.GetPopularTracks()(c)

		if err != nil {
			t2.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)

		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			t2.Fail()
		}
	})

}
