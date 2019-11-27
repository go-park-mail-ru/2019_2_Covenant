package middlewares

import (
	. "2019_2_Covenant/internal/models"
	mockSs "2019_2_Covenant/internal/session/mocks"
	mockUs "2019_2_Covenant/internal/user/mocks"
	"2019_2_Covenant/pkg/logger"
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

func TestMiddlewareManager_CheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)
	logrus.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T){
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		cookie := http.Cookie{Name: "Covenant", Value: "covenantcookies"}
		req.AddCookie(&cookie)

		c := e.NewContext(req, rec)

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		SUsecase.EXPECT().Get("covenantcookies").Return(sess, nil)

		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(sess.UserID).Return(user, nil)

		h := MiddlewareManager.CheckAuth(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err := h(c)

		if err != nil {
			fmt.Println("Error: expected nil, got", err)
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		expBody := `test`
		if strings.Trim(string(body), "\n") != expBody {
			fmt.Println("Body: expected:", expBody)
			fmt.Println("      got:", string(body))
			t1.Fail()
		}

		status := rec.Code
		if status != 200 {
			fmt.Println("Status: expected 200, got", status)
			t1.Fail()
		}
	})

	t.Run("Error unauthorized", func(t2 *testing.T){
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := MiddlewareManager.CheckAuth(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err := h(c)

		if err != nil {
			fmt.Println("Error: expected nil, got", err)
			t2.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		expBody := `{"error":"unauthorized"}`
		if strings.Trim(string(body), "\n") != expBody {
			fmt.Println("Body: expected:", expBody)
			fmt.Println("      got:", string(body))
			t2.Fail()
		}

		status := rec.Code
		if status != 401 {
			fmt.Println("Status: expected 401, got", status)
			t2.Fail()
		}
	})

	t.Run("Error unauthorized", func(t3 *testing.T){
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		cookie := http.Cookie{Name: "Covenant", Value: "covenantcookies"}
		req.AddCookie(&cookie)

		c := e.NewContext(req, rec)

		SUsecase.EXPECT().Get("covenantcookies").Return(nil, fmt.Errorf("some error"))

		h := MiddlewareManager.CheckAuth(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err := h(c)

		if err != nil {
			fmt.Println("Error: expected nil, got", err)
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		expBody := `{"error":"unauthorized"}`
		if strings.Trim(string(body), "\n") != expBody {
			fmt.Println("Body: expected:", expBody)
			fmt.Println("      got:", string(body))
			t3.Fail()
		}

		status := rec.Code
		if status != 401 {
			fmt.Println("Status: expected 401, got", status)
			t3.Fail()
		}
	})

	t.Run("Test OK", func(t4 *testing.T){
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		cookie := http.Cookie{Name: "Covenant", Value: "covenantcookies"}
		req.AddCookie(&cookie)

		c := e.NewContext(req, rec)

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		SUsecase.EXPECT().Get("covenantcookies").Return(sess, nil)

		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(sess.UserID).Return(nil, fmt.Errorf("some error"))

		h := MiddlewareManager.CheckAuth(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err := h(c)

		if err != nil {
			fmt.Println("Error: expected nil, got", err)
			t4.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		expBody := `{"error":"some error"}`
		if strings.Trim(string(body), "\n") != expBody {
			fmt.Println("Body: expected:", expBody)
			fmt.Println("      got:", string(body))
			t4.Fail()
		}

		status := rec.Code
		if status != 400 {
			fmt.Println("Status: expected 200, got", status)
			t4.Fail()
		}
	})
}

func TestMiddlewareManager_PanicRecovering(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)
	logrus.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		h := MiddlewareManager.PanicRecovering(func(c echo.Context) error {
			panic("some panic")
			return nil
		})

		err := h(c)

		if err != nil {
			fmt.Println("Error: expected nil, got", err)
			t1.Fail()
		}
	})
}

func TestMiddlewareManager_CORSMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)
	logrus.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T){
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.Request().Header.Set("Origin", "http://localhost:3000")

		h := MiddlewareManager.CORSMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err := h(c)

		if err != nil {
			fmt.Println("Error: expected nil, got", err)
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		expBody := `test`
		if strings.Trim(string(body), "\n") != expBody {
			fmt.Println("Body: expected:", expBody)
			fmt.Println("      got:", string(body))
			t1.Fail()
		}

		status := rec.Code
		if status != 200 {
			fmt.Println("Status: expected 200, got", status)
			t1.Fail()
		}
	})

	t.Run("Error mothod option", func(t1 *testing.T){
		e := echo.New()

		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.Request().Header.Set("Origin", "http://localhost:3000")

		h := MiddlewareManager.CORSMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err := h(c)

		if err != nil {
			fmt.Println("Error: expected nil, got", err)
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		expBody := ``
		if strings.Trim(string(body), "\n") != expBody {
			fmt.Println("Body: expected:", expBody)
			fmt.Println("      got:", string(body))
			t1.Fail()
		}
	})
}

func TestMiddlewareManager_CSRFCheckMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)
	logrus.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T){
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		token, _ :=  NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, sess.Expires)
		c.Request().Header.Set("X-Csrf-Token", token)

		h := MiddlewareManager.CSRFCheckMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err := h(c)

		if err != nil {
			fmt.Println("Error: expected nil, got", err)
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		expBody := `test`
		if strings.Trim(string(body), "\n") != expBody {
			fmt.Println("Body: expected:", expBody)
			fmt.Println("      got:", string(body))
			t1.Fail()
		}

		status := rec.Code
		if status != 200 {
			fmt.Println("Status: expected 200, got", status)
			t1.Fail()
		}
	})

	t.Run("Error with cookies", func(t2 *testing.T){
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		token, _ :=  NewCSRFTokenManager("Covenant").Create(sess.UserID, "wrong cookies", sess.Expires)
		c.Request().Header.Set("X-Csrf-Token", token)

		h := MiddlewareManager.CSRFCheckMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err := h(c)

		if err != nil {
			fmt.Println("Error: expected nil, got", err)
			t2.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		expBody := `{"error":"csrf error"}`
		if strings.Trim(string(body), "\n") != expBody {
			fmt.Println("Body: expected:", expBody)
			fmt.Println("      got:", string(body))
			t2.Fail()
		}

		status := rec.Code
		if status != 400 {
			fmt.Println("Status: expected 400, got", status)
			t2.Fail()
		}
	})

	t.Run("Error with token", func(t3 *testing.T){
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		c.Request().Header.Set("X-Csrf-Token", "wrong token")

		h := MiddlewareManager.CSRFCheckMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		err := h(c)

		if err != nil {
			fmt.Println("Error: expected nil, got", err)
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		expBody := `{"error":"csrf error"}`
		if strings.Trim(string(body), "\n") != expBody {
			fmt.Println("Body: expected:", expBody)
			fmt.Println("      got:", string(body))
			t3.Fail()
		}

		status := rec.Code
		if status != 500 {
			fmt.Println("Status: expected 500, got", status)
			t3.Fail()
		}
	})
}
