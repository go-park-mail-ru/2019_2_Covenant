package delivery

import (
	. "2019_2_Covenant/internal/middlewares"
	. "2019_2_Covenant/internal/models"
	mockSs "2019_2_Covenant/internal/session/mocks"
	mockUs "2019_2_Covenant/internal/user/mocks"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

//go:generate mockgen -source=../usecase.go -destination=../mocks/mock_usecase.go -package=mock

func TestSessionHandler_GetCSRF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	SUsecase := mockSs.NewMockRepository(ctrl)
	UUsecase := mockUs.NewMockRepository(ctrl)
	Logger := logrus.New()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewSessionHandler(SUsecase, MiddlewareManager, Logger)
	handler.Logger.Out = ioutil.Discard

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/get_csrf")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		err := handler.GetCSRF()(c)

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

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/get_csrf")

		err := handler.GetCSRF()(c)

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