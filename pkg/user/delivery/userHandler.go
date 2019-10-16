package delivery

import (
	"2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/user"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
)

type UserHandler struct {
	UUsecase user.Usecase
}

func NewUserHandler(e *echo.Echo, uUC user.Usecase) {
	handler := &UserHandler{
		UUsecase: uUC,
	}

	e.POST("/api/v1/signup", handler.SignUp)
	//e.POST("/api/v1/signin",
}

type ResponseError struct {
	Error string `json:"error"`
}

func isValidSignUpReq(usr models.UserReg) (bool, error) {
	v := validator.New()
	err := v.Struct(usr)

	if err != nil {
		return false, models.ErrBadParam
	}

	return true, nil
}

func isValidSignInReq(usr models.UserLogin) (bool, error) {
	v := validator.New()
	err := v.Struct(usr)

	if err != nil {
		return false, models.ErrBadParam
	}

	return true, nil
}

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
		fmt.Println("!")
		return c.JSON(http.StatusBadRequest, ResponseError{models.ErrAlreadyExist.Error()})
	}

	newUser := &models.User{
		Email: userRegData.Email,
		Password: userRegData.Password,
		Username: userRegData.Username,
	}

	err = uh.UUsecase.Store(newUser)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
	}

	cookie := &http.Cookie{
		Name:       "Covenant",
		Value:      uuid.New().String(),
		Expires:    time.Now().Add(24 * time.Hour),
	}

	session := &models.Session{
		UserID:  newUser.ID,
		Expires: cookie.Expires,
		Data:    cookie.Value,
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, newUser)
}

// curl -X POST 127.0.0.1:8000/api/v1/signin -H 'Content-Type: application/json' \
// curl -X POST 127.0.0.1:8000/api/v1/signup -H 'Content-Type: application/json' \
// -d '{"email": "m1@mail.ru", "username": "Marsha1l", "password": "12345312"}'
