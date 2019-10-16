package delivery

import (
	"2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/session"
	"2019_2_Covenant/pkg/user"
	"2019_2_Covenant/pkg/vars"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
)

type UserHandler struct {
	UUsecase     user.Usecase
	sesssionRepo session.Repository
}

func NewUserHandler(e *echo.Echo, uUC user.Usecase, sR session.Repository) {
	handler := &UserHandler{
		UUsecase:     uUC,
		sesssionRepo: sR,
	}

	e.POST("/api/v1/signup", handler.SignUp)

	e.POST("/api/v1/signin", handler.SignIn, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("Covenant")

			if err != nil {
				err = next(c)
				return err
			}

			sess, err := handler.sesssionRepo.Get(cookie.Value)

			if err != nil {
				err = next(c)
				return err
			}

			usr, err := handler.UUsecase.GetByID(sess.UserID)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}

			fmt.Println("authorized")
			return c.JSON(http.StatusOK, usr)
		}
	})
}

type ResponseError struct {
	Error string `json:"error"`
}

func isValidSignUpReq(usr models.UserReg) (bool, error) {
	v := validator.New()
	err := v.Struct(usr)

	if err != nil {
		return false, vars.ErrBadParam
	}

	return true, nil
}

func isValidSignInReq(usr models.UserLogin) (bool, error) {
	v := validator.New()
	err := v.Struct(usr)

	if err != nil {
		return false, vars.ErrBadParam
	}

	return true, nil
}

// curl -X POST 127.0.0.1:8000/api/v1/signup -H 'Content-Type: application/json' \
// -d '{"email": "m@mail.ru", "username": "Marshal", "password": "12345312"}'

func (uh UserHandler) SignUp(c echo.Context) error {
	var userRegData models.UserReg
	err := c.Bind(&userRegData)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{err.Error()})
	}

	if ok, err := isValidSignUpReq(userRegData); !ok {
		return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
	}

	usr, err := uh.UUsecase.GetByEmail(userRegData.Email)

	if usr != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{vars.ErrAlreadyExist.Error()})
	}

	newUser := &models.User{
		Email:    userRegData.Email,
		Password: userRegData.Password,
		Username: userRegData.Username,
	}

	err = uh.UUsecase.Store(newUser)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
	}

	cookie := &http.Cookie{
		Name:    "Covenant",
		Value:   uuid.New().String(),
		Expires: time.Now().Add(24 * time.Hour),
	}

	sess := &models.Session{
		UserID:  newUser.ID,
		Expires: cookie.Expires,
		Data:    cookie.Value,
	}

	err = uh.sesssionRepo.Store(sess)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, newUser)
}

// curl -X POST 127.0.0.1:8000/api/v1/signin -H 'Content-Type: application/json' \
// -d '{"email": "m@mail.ru", "password": "12345312"}'

func (uh UserHandler) SignIn(c echo.Context) error {
	var userLoginData models.UserLogin
	err := c.Bind(&userLoginData)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{err.Error()})
	}

	if ok, err := isValidSignInReq(userLoginData); !ok {
		return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
	}

	usr, err := uh.UUsecase.GetByEmail(userLoginData.Email)

	if usr == nil {
		return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
	}

	if usr.Password != userLoginData.Password {
		return c.JSON(http.StatusBadRequest, ResponseError{vars.ErrBadParam.Error()})
	}

	cookie := &http.Cookie{
		Name:    "Covenant",
		Value:   uuid.New().String(),
		Expires: time.Now().Add(24 * time.Hour),
	}

	sess := &models.Session{
		UserID:  usr.ID,
		Expires: cookie.Expires,
		Data:    cookie.Value,
	}

	err = uh.sesssionRepo.Store(sess)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, usr)
}
