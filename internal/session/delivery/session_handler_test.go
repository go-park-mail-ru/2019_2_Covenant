package delivery

import (
	. "2019_2_Covenant/internal/middlewares"
	. "2019_2_Covenant/internal/models"
	mockSs "2019_2_Covenant/internal/session/mocks"
	mockUs "2019_2_Covenant/internal/user/mocks"
	"2019_2_Covenant/pkg/logger"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

//go:generate mockgen -source=../usecase.go -destination=../mocks/mock_usecase.go -package=mock
//go:generate mockgen -source=../../user/usecase/usecase.go -destination=../../user/mocks/mock_usecase.go -package=mock

func TestSessionHandler_GetCSRF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	SUsecase := mockSs.NewMockRepository(ctrl)
	UUsecase := mockUs.NewMockRepository(ctrl)
	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewSessionHandler(SUsecase, UUsecase, MiddlewareManager, Logger)
	Logger.L.SetOutput(ioutil.Discard)

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

func TestSessionHandler_CreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	SUsecase := mockSs.NewMockRepository(ctrl)
	UUsecase := mockUs.NewMockRepository(ctrl)
	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewSessionHandler(SUsecase, UUsecase, MiddlewareManager, Logger)
	Logger.L.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		sessJSON := `{"email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(sessJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/session")

		usr := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = usr.BeforeStore()

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		UUsecase.EXPECT().GetByEmail("e@mail.ru").Return(usr, nil)
		SUsecase.EXPECT().Store(gomock.Any()).Return(nil)
		err := handler.CreateSession()(c)

		if err != nil {
			fmt.Println("Error happens")
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"body":{"user":{"nickname":"nickname","email":"e@mail.ru","avatar":"path","role":0,"access":0}}}` {
			fmt.Println(string(body))
			t1.Fail()
		}
	})

	t.Run("Error Bad params", func(t2 *testing.T) {
		e := echo.New()

		sessJSON := `{"email":"e@mail", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(sessJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/session")

		err := handler.CreateSession()(c)

		if err != nil {
			fmt.Println("Error happens")
			t2.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"bad params"}` {
			fmt.Println(string(body))
			t2.Fail()
		}
	})

	t.Run("Error getting by email", func(t3 *testing.T) {
		e := echo.New()

		sessJSON := `{"email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(sessJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/session")

		usr := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = usr.BeforeStore()

		UUsecase.EXPECT().GetByEmail("e@mail.ru").Return(nil, fmt.Errorf("some error"))
		err := handler.CreateSession()(c)

		if err != nil {
			fmt.Println("Error happens")
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"some error"}` {
			fmt.Println(string(body))
			t3.Fail()
		}
	})

	t.Run("Error verifying pass", func(t4 *testing.T) {
		e := echo.New()

		sessJSON := `{"email":"e@mail.ru", "password":"wrong pass"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(sessJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/session")

		usr := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = usr.BeforeStore()

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		UUsecase.EXPECT().GetByEmail("e@mail.ru").Return(usr, nil)
		err := handler.CreateSession()(c)

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

	t.Run("Error storing", func(t5 *testing.T) {
		e := echo.New()

		sessJSON := `{"email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(sessJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/session")

		usr := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = usr.BeforeStore()

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		UUsecase.EXPECT().GetByEmail("e@mail.ru").Return(usr, nil)
		SUsecase.EXPECT().Store(gomock.Any()).Return(fmt.Errorf("some error"))
		err := handler.CreateSession()(c)

		if err != nil {
			fmt.Println("Error happens")
			t5.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"some error"}` {
			fmt.Println(string(body))
			t5.Fail()
		}
	})
}

func TestSessionHandler_DeleteSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	SUsecase := mockSs.NewMockRepository(ctrl)
	UUsecase := mockUs.NewMockRepository(ctrl)
	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewSessionHandler(SUsecase, UUsecase, MiddlewareManager, Logger)
	Logger.L.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/session")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		SUsecase.EXPECT().DeleteByID(sess.ID).Return(nil)
		err := handler.DeleteSession()(c)

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

	t.Run("Error getting from ctx", func(t2 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/session")

		err := handler.DeleteSession()(c)

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

	t.Run("Error deleting", func(t3 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/session")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		SUsecase.EXPECT().DeleteByID(sess.ID).Return(fmt.Errorf("some error"))
		err := handler.DeleteSession()(c)

		if err != nil {
			fmt.Println("Error happens")
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"some error"}` {
			fmt.Println(string(body))
			t3.Fail()
		}
	})
}
