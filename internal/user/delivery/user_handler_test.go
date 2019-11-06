package delivery

import (
	. "2019_2_Covenant/internal/models"
	mockSs "2019_2_Covenant/internal/session/mocks"
	mockUs "2019_2_Covenant/internal/user/mocks"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"testing"
)

//go:generate mockgen -source=../usecase.go -destination=../mocks/mock_usecase.go -package=mock
//go:generate mockgen -source=../../session/usecase.go -destination=../../session/mocks/mock_usecase.go -package=mock

func TestUserHandler_LogIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	handler := UserHandler{UUsecase: UUsecase, SUsecase: SUsecase}

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		userJSON := `{"email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodGet, "/api/v1", strings.NewReader(userJSON))
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
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"body":{"nickname":"nick","email":"e@mail.ru","avatar":"path","role":0,"access":0}}` {
			t1.Fail()
		}
	})

	t.Run("Error of validating", func(t2 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/login")
		err := handler.LogIn()(c)

		if err != nil {
			t2.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"bad params"}` {
			t2.Fail()
		}
	})

	t.Run("Error of verifying", func(t3 *testing.T) {
		e := echo.New()

		userJSON := `{"email":"e@mail.ru", "password":"not real"}`
		req := httptest.NewRequest(http.MethodGet, "/api/v1", strings.NewReader(userJSON))
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
		req := httptest.NewRequest(http.MethodGet, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/login")

		email := "e@mail.ru"

		UUsecase.EXPECT().GetByEmail(email).Return(nil, fmt.Errorf("some err"))

		err := handler.LogIn()(c)
		if err != nil {
			t4.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"some err"}` {
			t4.Fail()
		}
	})

	t.Run("Error storing session", func(t5 *testing.T) {
		e := echo.New()

		userJSON := `{"email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodGet, "/api/v1", strings.NewReader(userJSON))
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
			t5.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			t5.Fail()
		}
	})
}

func TestUserHandler_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	handler := UserHandler{UUsecase: UUsecase, SUsecase: SUsecase}

	t.Run("Test OK", func(t1 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"nickname","email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodGet, "/api/v1", strings.NewReader(userJSON))
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
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		//fmt.Println(string(body))
		if strings.Trim(string(body), "\n") != `{"body":{"nickname":"nickname","email":"e@mail.ru","avatar":"","role":0,"access":0}}` {
			t1.Fail()
		}
	})

	t.Run("Error of validating", func(t2 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/signup")
		err := handler.SignUp()(c)

		if err != nil {
			t2.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		if strings.Trim(string(body), "\n") != `{"error":"bad params"}` {
			t2.Fail()
		}
	})

	t.Run("Error already exist", func(t3 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"nickname","email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodGet, "/api/v1", strings.NewReader(userJSON))
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
			t3.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		//fmt.Println(string(body))
		if strings.Trim(string(body), "\n") != `{"error":"already exist"}` {
			t3.Fail()
		}
	})

	t.Run("Error storing user", func(t4 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"nickname","email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodGet, "/api/v1", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/signup")

		email := "e@mail.ru"

		UUsecase.EXPECT().GetByEmail(email).Return(nil, fmt.Errorf("some error"))

		UUsecase.EXPECT().Store(gomock.Any()).Return(nil, fmt.Errorf("some error"))
		err := handler.SignUp()(c)

		if err != nil {
			t4.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		//fmt.Println(string(body))
		if strings.Trim(string(body), "\n") != `{"error":"some error"}` {
			t4.Fail()
		}
	})

	t.Run("Error storing session", func(t5 *testing.T) {
		e := echo.New()

		userJSON := `{"nickname":"nickname","email":"e@mail.ru", "password":"qwerty"}`
		req := httptest.NewRequest(http.MethodGet, "/api/v1", strings.NewReader(userJSON))
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
			t5.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		//fmt.Println(string(body))
		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			t5.Fail()
		}
	})
}

func TestUserHandler_EditProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UUsecase := mockUs.NewMockRepository(ctrl)
	SUsecase := mockSs.NewMockRepository(ctrl)

	handler := UserHandler{UUsecase: UUsecase, SUsecase: SUsecase}

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
			Expires: time.Now().Add(24*time.Hour),
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
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		//fmt.Println(string(body))
		if strings.Trim(string(body), "\n") != `{"body":{"nickname":"new_nickname","email":"e@mail.ru","avatar":"path","role":0,"access":0}}` {
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
			t2.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		//fmt.Println(string(body))
		if strings.Trim(string(body), "\n") != `{"error":"internal server error"}` {
			t2.Fail()
		}
	})

	t.Run("Error bad params", func(t1 *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/api/v1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		sess := &Session{
			ID:      1,
			UserID:  2,
			Expires: time.Now().Add(24*time.Hour),
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
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		//fmt.Println(string(body))
		if strings.Trim(string(body), "\n") != `{"error":"bad params"}` {
			t1.Fail()
		}
	})

	t.Run("Error updating", func(t1 *testing.T) {
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
			Expires: time.Now().Add(24*time.Hour),
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
			t1.Fail()
		}

		body, _ := ioutil.ReadAll(rec.Body)
		//fmt.Println(string(body))
		if strings.Trim(string(body), "\n") != `{"error":"some error"}` {
			t1.Fail()
		}
	})
}
