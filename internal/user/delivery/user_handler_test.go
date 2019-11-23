package delivery

import (
	. "2019_2_Covenant/internal/middlewares"
	. "2019_2_Covenant/internal/models"
	mockSs "2019_2_Covenant/internal/session/mocks"
	mockUs "2019_2_Covenant/internal/user/mocks"
	"2019_2_Covenant/pkg/logger"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

//TODO: уточнить у Марселя по поводу валидации. Сейчас не протестировано

//go:generate mockgen -source=../usecase.go -destination=../mocks/mock_usecase.go -package=mock
//go:generate mockgen -source=../../session/usecase.go -destination=../../session/mocks/mock_usecase.go -package=mock

func TestUserHandler_LogIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewUserHandler(UUsecase, SUsecase, MiddlewareManager, Logger)
	logrus.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		userJSON := `{"email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/login")

		email := "e@mail.ru"
		user := &User{
			ID: 1, Nickname: "nick", Email: email, PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByEmail(email).Return(user, nil)

		SUsecase.EXPECT().Store(gomock.Any()).Return(nil)
		err := handler.LogIn()(c)

		if err != nil {
			fmt.Println("Error happens")
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"body":{"nickname":"nick","email":"e@mail.ru","avatar":"path","role":0,"access":0}}` {
			fmt.Println(string(body))
			t1.Fail()
		}
	})

	t.Run("Error of validating", func(t2 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/login")
		err := handler.LogIn()(c)

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

	t.Run("Error of verifying", func(t3 *testing.T) {
		e := echo.New()

		userJSON := `{"email":"e@mail.ru", "password":"not real"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/login")

		email := "e@mail.ru"
		user := &User{
			ID: 1, Nickname: "nick", Email: email, PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByEmail(email).Return(user, nil)

		err := handler.LogIn()(c)
		if err != nil {
			fmt.Println("Error happens")
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"bad params"}` {
			t3.Fail()
		}
	})

	t.Run("Error not exist", func(t4 *testing.T) {
		e := echo.New()

		userJSON := `{"email":"e@mail.ru", "password":"not real"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/login")

		email := "e@mail.ru"

		UUsecase.EXPECT().GetByEmail(email).Return(nil, fmt.Errorf("some err"))

		err := handler.LogIn()(c)
		if err != nil {
			fmt.Println("Error happens")
			t4.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"some err"}` {
			fmt.Println(string(body))
			t4.Fail()
		}
	})

	t.Run("Error storing session", func(t5 *testing.T) {
		e := echo.New()

		userJSON := `{"email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/login")

		email := "e@mail.ru"
		user := &User{
			ID: 1, Nickname: "nick", Email: email, PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByEmail(email).Return(user, nil)

		SUsecase.EXPECT().Store(gomock.Any()).Return(fmt.Errorf("some error"))
		err := handler.LogIn()(c)

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

	t.Run("Error of binding", func(t6 *testing.T) {
		e := echo.New()

		userJSON := `{"mistake"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/login")

		err := handler.LogIn()(c)

		if err != nil {
			fmt.Println("Error happens", err)
			t6.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"code=400, message=Syntax error: offset=11, error=invalid character '}' after object key, internal=invalid character '}' after object key"}` {
			fmt.Println(string(body))
			t6.Fail()
		}
	})
}

func TestUserHandler_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewUserHandler(UUsecase, SUsecase, MiddlewareManager, Logger)
	logrus.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"nickname","email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/signup")

		email := "e@mail.ru"

		UUsecase.EXPECT().GetByEmail(email).Return(nil, fmt.Errorf("some error"))

		user := &User{
			ID: 2, Nickname: "nickname", Email: email, PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()

		UUsecase.EXPECT().Store(gomock.Any()).Return(user, nil)
		SUsecase.EXPECT().Store(gomock.Any()).Return(nil)
		err := handler.SignUp()(c)

		if err != nil {
			fmt.Println("Error happens")
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"body":{"nickname":"nickname","email":"e@mail.ru","avatar":"","role":0,"access":0}}` {
			fmt.Println(string(body))
			t1.Fail()
		}
	})

	t.Run("Error of validating", func(t2 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/signup")
		err := handler.SignUp()(c)

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

	t.Run("Error already exist", func(t3 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"nickname","email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/signup")

		email := "e@mail.ru"

		user := &User{
			ID: 2, Nickname: "nickname", Email: email, PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()

		UUsecase.EXPECT().GetByEmail(email).Return(user, nil)

		err := handler.SignUp()(c)

		if err != nil {
			fmt.Println("Error happens")
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"already exist"}` {
			fmt.Println(string(body))
			t3.Fail()
		}
	})

	t.Run("Error storing user", func(t4 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"nickname","email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/signup")

		email := "e@mail.ru"

		UUsecase.EXPECT().GetByEmail(email).Return(nil, fmt.Errorf("some error"))

		UUsecase.EXPECT().Store(gomock.Any()).Return(nil, fmt.Errorf("some error"))
		err := handler.SignUp()(c)

		if err != nil {
			fmt.Println("Error happens")
			t4.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"some error"}` {
			fmt.Println(string(body))
			t4.Fail()
		}
	})

	t.Run("Error storing session", func(t5 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"nickname","email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/signup")

		email := "e@mail.ru"

		UUsecase.EXPECT().GetByEmail(email).Return(nil, fmt.Errorf("some error"))

		user := &User{
			ID: 2, Nickname: "nickname", Email: email, PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()

		UUsecase.EXPECT().Store(gomock.Any()).Return(user, nil)
		SUsecase.EXPECT().Store(gomock.Any()).Return(fmt.Errorf("some error"))
		err := handler.SignUp()(c)

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

	t.Run("Error of binding", func(t6 *testing.T) {
		e := echo.New()

		userJSON := `{"mistake"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/login")

		err := handler.SignUp()(c)

		if err != nil {
			fmt.Println("Error happens", err)
			t6.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"code=400, message=Syntax error: offset=11, error=invalid character '}' after object key, internal=invalid character '}' after object key"}` {
			fmt.Println(string(body))
			t6.Fail()
		}
	})
}

func TestUserHandler_EditProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewUserHandler(UUsecase, SUsecase, MiddlewareManager, Logger)
	logrus.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"new_nickname"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)
		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(uint64(2)).Return(user, nil)

		user.Nickname = "new_nickname"
		UUsecase.EXPECT().UpdateNickname(uint64(2), "new_nickname").Return(user, nil)
		err := handler.EditProfile()(c)

		if err != nil {
			fmt.Println("Error happens")
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"body":{"nickname":"new_nickname","email":"e@mail.ru","avatar":"path","role":0,"access":0}}` {
			fmt.Println(string(body))
			t1.Fail()
		}
	})

	t.Run("Error getting from ctx", func(t2 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"new_nickname"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		err := handler.EditProfile()(c)

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

	t.Run("Error bad params", func(t3 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)
		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(uint64(2)).Return(user, nil)

		err := handler.EditProfile()(c)

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

	t.Run("Error updating", func(t4 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"new_nickname"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)
		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(uint64(2)).Return(user, nil)

		user.Nickname = "new_nickname"
		UUsecase.EXPECT().UpdateNickname(uint64(2), "new_nickname").Return(nil, fmt.Errorf("some error"))
		err := handler.EditProfile()(c)

		if err != nil {
			fmt.Println("Error happens")
			t4.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"some error"}` {
			fmt.Println(string(body))
			t4.Fail()
		}
	})
	t.Run("Error getting by id", func(t5 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"new_nickname"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)
		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(uint64(2)).Return(nil, fmt.Errorf("some error"))

		err := handler.EditProfile()(c)

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

	t.Run("Error of binding", func(t6 *testing.T) {
		e := echo.New()

		userJSON := `{"mistake"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/login")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(uint64(2)).Return(user, nil)

		err := handler.EditProfile()(c)

		if err != nil {
			fmt.Println("Error happens", err)
			t6.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"code=400, message=Syntax error: offset=11, error=invalid character '}' after object key, internal=invalid character '}' after object key"}` {
			fmt.Println(string(body))
			t6.Fail()
		}
	})
}

func TestUserHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewUserHandler(UUsecase, SUsecase, MiddlewareManager, Logger)
	logrus.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(uint64(2)).Return(user, nil)

		err := handler.GetProfile()(c)

		if err != nil {
			fmt.Println("Error happens")
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"body":{"nickname":"nickname","email":"e@mail.ru","avatar":"path","role":0,"access":0}}` {
			fmt.Println(string(body))
			t1.Fail()
		}
	})

	t.Run("Error getting from ctx", func(t2 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		err := handler.GetProfile()(c)

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

	t.Run("Error getting by id", func(t3 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)
		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "path", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(uint64(2)).Return(nil, fmt.Errorf("some error"))

		err := handler.GetProfile()(c)

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

func TestUserHandler_LogOut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewUserHandler(UUsecase, SUsecase, MiddlewareManager, Logger)
	logrus.SetOutput(ioutil.Discard)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/logout")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		SUsecase.EXPECT().DeleteByID(sess.ID).Return(nil)

		err := handler.LogOut()(c)

		if err != nil {
			fmt.Println("Error happens")
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"message":"logout"}` {
			fmt.Println(string(body))
			t1.Fail()
		}
	})

	t.Run("Error getting from ctx", func(t2 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/logout")

		err := handler.LogOut()(c)

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

	t.Run("Error getting by id", func(t3 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/logout")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		SUsecase.EXPECT().DeleteByID(sess.ID).Return(fmt.Errorf("some error"))

		err := handler.LogOut()(c)

		if err != nil {
			fmt.Println("Error happens")
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			fmt.Println(string(body))
			t3.Fail()
		}
	})
}

func TestUserHandler_GetAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewUserHandler(UUsecase, SUsecase, MiddlewareManager, Logger)
	logrus.SetOutput(ioutil.Discard)

	rootPath := "/Users/yulia_plaksina/back/2019_2_Covenant"
	_ = os.Chdir(rootPath)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/avatar")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "resources/avatars/image.png", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(uint64(2)).Return(user, nil)

		err := handler.GetAvatar()(c)

		if err != nil {
			fmt.Println("Error happens", err)
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if body == nil {
			fmt.Println("Expected avatar, got:", string(body))
			t1.Fail()
		}
	})

	t.Run("Error getting from ctx", func(t2 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/avatar")

		err := handler.GetAvatar()(c)

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

	t.Run("Error getting by id", func(t3 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/avatar")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)
		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru", PlainPassword: "qwerty", Avatar: "image.png", Role: 0, Access: 0,
		}
		_ = user.BeforeStore()
		UUsecase.EXPECT().GetByID(uint64(2)).Return(nil, fmt.Errorf("some error"))

		err := handler.GetAvatar()(c)

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

func TestUserHandler_SetAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	Logger := logger.NewLogrusLogger()
	MiddlewareManager := NewMiddlewareManager(UUsecase, SUsecase, Logger)

	handler := NewUserHandler(UUsecase, SUsecase, MiddlewareManager, Logger)
	logrus.SetOutput(ioutil.Discard)

	rootPath := "/Users/yulia_plaksina/back/2019_2_Covenant"
	_ = os.Chdir(rootPath)

	avatarsPath := "/resources/avatars/image.png"
	filePath := filepath.Join(rootPath, avatarsPath)

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error of opening")
			t1.Fail()
		}
		defer file.Close()

		reqBody := &bytes.Buffer{}
		writer := multipart.NewWriter(reqBody)

		part, err := writer.CreateFormFile("avatar", filepath.Base(filePath))
		if err != nil {
			fmt.Println("Error of opening")
			t1.Fail()
		}
		_, err = io.Copy(part, file)

		err = writer.Close()
		if err != nil {
			fmt.Println("Error of closing")
			t1.Fail()
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1", reqBody)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.SetPath("/avatar")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		user := &User{
			ID: 2, Nickname: "nickname", Email: "e@mail.ru",
		}

		UUsecase.EXPECT().UpdateAvatar(uint64(2), gomock.Any()).Return(user, nil)
		err = handler.SetAvatar()(c)

		if err != nil {
			fmt.Println("Error happens", err)
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"body":{"nickname":"nickname","email":"e@mail.ru","avatar":"","role":0,"access":0}}` {
			fmt.Println(string(body))
			t1.Fail()
		}
	})

	t.Run("Error extracting files", func(t3 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/avatar")

		_ = c.File("not a file")

		err := handler.SetAvatar()(c)

		if err != nil {
			fmt.Println("Error happens", err)
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"retrieving error"}` {
			fmt.Println(string(body))
			t3.Fail()
		}
	})

	t.Run("Error session", func(t3 *testing.T) {
		e := echo.New()

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error of opening")
			t3.Fail()
		}
		defer file.Close()

		reqBody := &bytes.Buffer{}
		writer := multipart.NewWriter(reqBody)

		part, err := writer.CreateFormFile("avatar", filepath.Base(filePath))
		if err != nil {
			fmt.Println("Error of opening")
			t3.Fail()
		}
		_, err = io.Copy(part, file)

		err = writer.Close()
		if err != nil {
			fmt.Println("Error of closing")
			t3.Fail()
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1", reqBody)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.SetPath("/avatar")

		err = handler.SetAvatar()(c)

		if err != nil {
			fmt.Println("Error happens", err)
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			fmt.Println(string(body))
			t3.Fail()
		}
	})

	t.Run("Error with directory", func(t4 *testing.T) {

		e := echo.New()

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error of opening")
			t4.Fail()
		}
		defer file.Close()

		reqBody := &bytes.Buffer{}
		writer := multipart.NewWriter(reqBody)

		part, err := writer.CreateFormFile("avatar", filepath.Base(filePath))
		if err != nil {
			fmt.Println("Error of opening")
			t4.Fail()
		}
		_, err = io.Copy(part, file)

		err = writer.Close()
		if err != nil {
			fmt.Println("Error of closing")
			t4.Fail()
		}
		_ = os.Chdir("/Users/yulia_plaksina/back")

		req := httptest.NewRequest(http.MethodPost, "/api/v1", reqBody)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.SetPath("/avatar")

		err = handler.SetAvatar()(c)

		_ = os.Chdir(rootPath)
		if err != nil {
			fmt.Println("Error happens", err)
			t4.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			fmt.Println(string(body))
			t4.Fail()
		}
	})

	t.Run("Error updating", func(t5 *testing.T) {
		e := echo.New()

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error of opening")
			t5.Fail()
		}
		defer file.Close()

		reqBody := &bytes.Buffer{}
		writer := multipart.NewWriter(reqBody)

		part, err := writer.CreateFormFile("avatar", filepath.Base(filePath))
		if err != nil {
			fmt.Println("Error of opening")
			t5.Fail()
		}
		_, err = io.Copy(part, file)

		err = writer.Close()
		if err != nil {
			fmt.Println("Error of closing")
			t5.Fail()
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1", reqBody)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.SetPath("/avatar")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24 * time.Hour),
			Data:    "covenantcookies",
		}
		c.Set("session", sess)

		UUsecase.EXPECT().UpdateAvatar(uint64(2), gomock.Any()).Return(nil, fmt.Errorf("some error"))
		err = handler.SetAvatar()(c)

		if err != nil {
			fmt.Println("Error happens", err)
			t5.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			fmt.Println(string(body))
			t5.Fail()
		}
	})
}
