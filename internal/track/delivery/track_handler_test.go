package delivery

import (
	. "2019_2_Covenant/internal/middlewares"
	. "2019_2_Covenant/internal/models"
	mockSs "2019_2_Covenant/internal/session/mocks"
	mock "2019_2_Covenant/internal/track/mocks"
	mockUs "2019_2_Covenant/internal/user/mocks"
	"2019_2_Covenant/pkg/logger"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

// Для тестирования только этого файла:
// go test -v -cover -race ./internal/track/delivery

//go:generate mockgen -source=../usecase.go -destination=../mocks/mock_usecase.go -package=mock
//go:generate mockgen -source=../../session/usecase.go -destination=../../session/mocks/mock_usecase.go -package=mock

func TestTrackHandler_GetPopularTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TUsecase := mock.NewMockUsecase(ctrl)

	SUsecase := mockSs.NewMockRepository(ctrl)
	UUsecase := mockUs.NewMockRepository(ctrl)
	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewTrackHandler(TUsecase, MiddlewareManager, Logger)
	Logger.L.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/popular")

		tracks := []*Track{
			{ID: 1, Name: "Still loving you", Duration: "2019-10-31T00:06:28Z"},
		}

		TUsecase.EXPECT().FetchPopular(gomock.Any(), gomock.Any()).Return(tracks, nil)
		err := handler.GetPopularTracks()(c)

		if err != nil {
			fmt.Println("Error happens")
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)

		if strings.Trim(string(body), "\n") != `{"body":{"tracks":[{"id":1,"name":"Still loving you","duration":"00:06:28","photo":"","artist":"","album":"","path":""}]}}` {
			fmt.Println(string(body))
			t1.Fail()
		}
	})

	t.Run("Error", func(t2 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/popular")

		TUsecase.EXPECT().FetchPopular(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("some error"))
		err := handler.GetPopularTracks()(c)

		if err != nil {
			fmt.Println("Error happens")
			t2.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)

		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			fmt.Println(string(body))
			t2.Fail()
		}
	})
}

func TestTrackHandler_AddToFavourites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TUsecase := mock.NewMockUsecase(ctrl)

	SUsecase := mockSs.NewMockRepository(ctrl)
	UUsecase := mockUs.NewMockRepository(ctrl)
	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewTrackHandler(TUsecase, MiddlewareManager, Logger)
	Logger.L.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		dataToAdd := `{"track_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(dataToAdd))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		sess := &Session{
			ID:      uint64(1),
			UserID:  uint64(2),
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		TUsecase.EXPECT().StoreFavourite(sess.UserID, uint64(1)).Return(nil)
		err := handler.AddToFavourites()(c)

		if err != nil {
			fmt.Println("Error happens")
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)

		if strings.Trim(string(body), "\n") != `{"message":"success"}` {
			fmt.Println(string(body))
			t1.Fail()
		}
	})

	t.Run("Error getting session", func(t2 *testing.T) {
		e := echo.New()

		dataToAdd := `{"track_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(dataToAdd))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		err := handler.AddToFavourites()(c)

		if err != nil {
			fmt.Println("Error happens")
			t2.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		//fmt.Println(string(body))
		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			t2.Fail()
		}
	})

	t.Run("Error empty body", func(t3 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		sess := &Session{
			ID:      uint64(1),
			UserID:  uint64(2),
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		err := handler.AddToFavourites()(c)

		if err != nil {
			fmt.Println("Error happens")
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)

		if strings.Trim(string(body), "\n") != `{"error":"bad params"}` {
			fmt.Println(string(body))
			t3.Fail()
		}
	})

	t.Run("Error bad params", func(t4 *testing.T) {
		e := echo.New()

		dataToAdd := `{"other_id": 1}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(dataToAdd))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		sess := &Session{
			ID:      uint64(1),
			UserID:  uint64(2),
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		err := handler.AddToFavourites()(c)

		if err != nil {
			fmt.Println("Error happens")
			t4.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"bad params"}` {
			fmt.Println(string(body))
			t4.Fail()
		}
	})

	t.Run("Error already added", func(t5 *testing.T) {
		e := echo.New()

		dataToAdd := `{"track_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(dataToAdd))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		sess := &Session{
			ID:      uint64(1),
			UserID:  uint64(2),
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		TUsecase.EXPECT().StoreFavourite(sess.UserID, uint64(1)).Return(fmt.Errorf("already exist"))
		err := handler.AddToFavourites()(c)

		if err != nil {
			fmt.Println("Error happens")
			t5.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)

		if strings.Trim(string(body), "\n") != `{"error":"already exist"}` {
			fmt.Println(string(body))
			t5.Fail()
		}
	})

	t.Run("Error binding", func(t6 *testing.T) {
		e := echo.New()

		dataToAdd := `{"other_id"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(dataToAdd))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		sess := &Session{
			ID:      uint64(1),
			UserID:  uint64(2),
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		err := handler.AddToFavourites()(c)

		if err != nil {
			fmt.Println("Error happens")
			t6.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"unprocessable entity"}` {
			fmt.Println(string(body))
			t6.Fail()
		}
	})
}

func TestTrackHandler_RemoveFavourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TUsecase := mock.NewMockUsecase(ctrl)

	SUsecase := mockSs.NewMockRepository(ctrl)
	UUsecase := mockUs.NewMockRepository(ctrl)
	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewTrackHandler(TUsecase, MiddlewareManager, Logger)
	Logger.L.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		dataToRemove := `{"track_id":1}`
		req := httptest.NewRequest(http.MethodDelete, "/api/v1", strings.NewReader(dataToRemove))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		sess := &Session{
			ID:      uint64(1),
			UserID:  uint64(2),
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		TUsecase.EXPECT().RemoveFavourite(sess.UserID, uint64(1)).Return(nil)
		err := handler.RemoveFavourite()(c)

		if err != nil {
			fmt.Println("Error happens")
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)

		if strings.Trim(string(body), "\n") != `{"message":"success"}` {
			fmt.Println(string(body))
			t1.Fail()
		}
	})

	t.Run("Error getting session", func(t2 *testing.T) {
		e := echo.New()

		dataToRemove := `{"track_id":1}`
		req := httptest.NewRequest(http.MethodDelete, "/api/v1", strings.NewReader(dataToRemove))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		err := handler.RemoveFavourite()(c)

		if err != nil {
			fmt.Println("Error happens")
			t2.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)

		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			fmt.Println(string(body))
			t2.Fail()
		}
	})

	t.Run("Error empty body", func(t3 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodDelete, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		sess := &Session{
			ID:      uint64(1),
			UserID:  uint64(2),
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		err := handler.RemoveFavourite()(c)

		if err != nil {
			fmt.Println("Error happens")
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)

		if strings.Trim(string(body), "\n") != `{"error":"bad params"}` {
			fmt.Println(string(body))
			t3.Fail()
		}
	})

	t.Run("Error bad params", func(t4 *testing.T) {
		e := echo.New()

		dataToRemove := `{"other_id":1}`
		req := httptest.NewRequest(http.MethodDelete, "/api/v1", strings.NewReader(dataToRemove))
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		sess := &Session{
			ID:      uint64(1),
			UserID:  uint64(2),
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		err := handler.RemoveFavourite()(c)

		if err != nil {
			fmt.Println("Error happens")
			t4.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"unprocessable entity"}` {
			fmt.Println(string(body))
			t4.Fail()
		}
	})

	t.Run("Error not found", func(t5 *testing.T) {
		e := echo.New()

		dataToRemove := `{"track_id":1}`
		req := httptest.NewRequest(http.MethodDelete, "/api/v1", strings.NewReader(dataToRemove))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		sess := &Session{
			ID:      uint64(1),
			UserID:  uint64(2),
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		TUsecase.EXPECT().RemoveFavourite(sess.UserID, uint64(1)).Return(fmt.Errorf("not found"))
		err := handler.RemoveFavourite()(c)

		if err != nil {
			fmt.Println("Error happens")
			t5.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)

		if strings.Trim(string(body), "\n") != `{"error":"not found"}` {
			fmt.Println(string(body))
			t5.Fail()
		}
	})

	t.Run("Error binding", func(t6 *testing.T) {
		e := echo.New()

		dataToAdd := `{"other_id"}`
		req := httptest.NewRequest(http.MethodDelete, "/api/v1", strings.NewReader(dataToAdd))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/tracks/favourite")

		sess := &Session{
			ID:      uint64(1),
			UserID:  uint64(2),
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		err := handler.RemoveFavourite()(c)

		if err != nil {
			fmt.Println("Error happens")
			t6.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"unprocessable entity"}` {
			fmt.Println(string(body))
			t6.Fail()
		}
	})
}

//func TestTrackHandler_GetFavourites(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	TUsecase := mock.NewMockUsecase(ctrl)
//
//	SUsecase := mockSs.NewMockRepository(ctrl)
//	UUsecase := mockUs.NewMockRepository(ctrl)
//	Logger := logger.NewLogrusLogger()
//	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)
//
//	handler := NewTrackHandler(TUsecase, MiddlewareManager, Logger)
//	Logger.L.SetOutput(ioutil.Discard)
//
//	t.Run("Test OK", func(t1 *testing.T) {
//		e := echo.New()
//
//		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
//		rec := httptest.NewRecorder()
//
//		c := e.NewContext(req, rec)
//		c.SetPath("/tracks/favourite")
//
//		sess := &Session{
//			ID:      uint64(1),
//			UserID:  uint64(2),
//			Expires: time.Now().Add(24 * time.Hour),
//			Data:    "covenantcookies",
//		}
//		c.Set("session", sess)
//
//		tracks := []*Track{
//			{ID: 1, Name: "Still loving you", Duration: "2019-10-31T00:06:28Z"},
//		}
//
//		TUsecase.EXPECT().FetchFavourites(sess.UserID, gomock.Any(), gomock.Any()).Return(tracks, uint64(1), nil)
//		err := handler.GetFavourites()(c)
//
//		if err != nil {
//			fmt.Println("Error happens")
//			t1.Fail()
//		}
//
//		body, _ := ioutil.ReadAll(rec.Body)
//
//		if strings.Trim(string(body), "\n") != `{"body":[{"id":1,"name":"Still loving you","duration":"00:06:28","photo":"","artist":"","album":"","path":""}]}` {
//			fmt.Println(string(body))
//			t1.Fail()
//		}
//	})
//
//	t.Run("Error getting session", func(t2 *testing.T) {
//		e := echo.New()
//
//		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
//		rec := httptest.NewRecorder()
//
//		c := e.NewContext(req, rec)
//		c.SetPath("/tracks/favourite")
//
//		err := handler.GetFavourites()(c)
//
//		if err != nil {
//			fmt.Println("Error happens")
//			t2.Fail()
//		}
//
//		body, _ := ioutil.ReadAll(rec.Body)
//
//		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
//			fmt.Println(string(body))
//			t2.Fail()
//		}
//	})
//
//	t.Run("Error with db", func(t3 *testing.T) {
//		e := echo.New()
//
//		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
//		rec := httptest.NewRecorder()
//
//		c := e.NewContext(req, rec)
//		c.SetPath("/tracks/favourite")
//
//		sess := &Session{
//			ID:      uint64(1),
//			UserID:  uint64(2),
//			Expires: time.Now().Add(24 * time.Hour),
//			Data:    "covenantcookies",
//		}
//		c.Set("session", sess)
//
//		TUsecase.EXPECT().FetchFavourites(sess.UserID, gomock.Any(), gomock.Any()).Return(nil, uint64(3), fmt.Errorf("some error"))
//		err := handler.GetFavourites()(c)
//
//		if err != nil {
//			fmt.Println("Error happens")
//			t3.Fail()
//		}
//
//		body, _ := ioutil.ReadAll(rec.Body)
//
//		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
//			fmt.Println(string(body))
//			t3.Fail()
//		}
//	})
//}
