package delivery

import (
	"2019_2_Covenant/internal/middleware"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/internal/vars"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UserHandler struct {
	UUsecase user.Usecase
	SUsecase session.Usecase
	MManager middleware.MiddlewareManager
}

func NewUserHandler(uUC user.Usecase, sUC session.Usecase) *UserHandler {
	return &UserHandler{
		UUsecase: uUC,
		SUsecase: sUC,
		MManager: middleware.NewMiddlewareManager(uUC, sUC),
	}
}

func (uh *UserHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/signup", uh.SignUp())
	e.POST("/api/v1/signin", uh.SignIn())
	e.POST("/api/v1/profile", uh.EditProfile(), uh.MManager.CheckAuth)
	e.GET("/api/v1/profile", uh.GetProfile, uh.MManager.CheckAuth)
	e.POST("/api/v1/avatar", uh.SetAvatar, uh.MManager.CheckAuth)
	e.GET("/api/v1/avatar", uh.GetAvatar, uh.MManager.CheckAuth)
}

type ResponseError struct {
	Error string `json:"error"`
}

type Response struct {
	Body interface{} `json:"body"`
}

func isValidRequest(usr interface{}) (bool, error) {
	v := validator.New()
	err := v.Struct(usr)

	if err != nil {
		return false, vars.ErrBadParam
	}

	return true, nil
}

// @Summary SignUp Route
// @Description Signing user up
// @ID sign-up-user
// @Accept json
// @Produce json
// @Param Data body models.UserReg true "JSON that contains user sign up data"
// @Success 200 object models.User
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/signup [post]
func (uh *UserHandler) SignUp() echo.HandlerFunc {
	type UserReg struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,gte=6"`
	}

	return func(c echo.Context) error {
		var userRegData UserReg
		err := c.Bind(&userRegData)

		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, ResponseError{err.Error()})
		}

		if ok, err := isValidRequest(userRegData); !ok {
			return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
		}

		usr, err := uh.UUsecase.GetByEmail(userRegData.Email)

		if usr != nil {
			return c.JSON(http.StatusBadRequest, ResponseError{vars.ErrAlreadyExist.Error()})
		}

		newUser := &models.User{
			Email:    userRegData.Email,
			Password: userRegData.Password,
			Nickname: userRegData.Username,
		}

		user, err := uh.UUsecase.Store(newUser)

		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
		}

		cookie := &http.Cookie{
			Name:    "Covenant",
			Value:   uuid.New().String(),
			Expires: time.Now().Add(24 * time.Hour),
		}

		sess := &models.Session{
			UserID:  user.ID,
			Expires: cookie.Expires,
			Data:    cookie.Value,
		}

		err = uh.SUsecase.Store(sess)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, Response{newUser})
	}
}


// @Summary SignIn Route
// @Description Signing user in
// @ID sign-in-user
// @Accept json
// @Produce json
// @Param Data body models.UserLogin true "JSON that contains user login data"
// @Success 200 object models.User
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/signin [post]
func (uh *UserHandler) SignIn() echo.HandlerFunc {
	type UserLogin struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	return func(c echo.Context) error {
		var userLoginData UserLogin
		err := c.Bind(&userLoginData)

		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, ResponseError{err.Error()})
		}

		if ok, err := isValidRequest(userLoginData); !ok {
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

		err = uh.SUsecase.Store(sess)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, Response{usr})
	}
}

// @Summary Edit Profile Route
// @Description Signing user in
// @ID edit-profile
// @Accept json
// @Produce json
// @Param Data body models.UserEdit true "JSON that contains user data to edit"
// @Success 200 object models.User
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/profile [post]
func (uh *UserHandler) EditProfile() echo.HandlerFunc {
	type UserEdit struct {
		Name    string `json:"name" validate:"required"`
		Surname string `json:"surname" validate:"required"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		usr, err := uh.UUsecase.GetByID(sess.UserID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
		}

		var userEditData UserEdit
		err = c.Bind(&userEditData)

		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, ResponseError{err.Error()})
		}

		if ok, err := isValidRequest(userEditData); !ok {
			return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
		}

		if usr, err = uh.UUsecase.Update(usr.ID, userEditData.Name, userEditData.Surname); err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseError{err.Error()})
		}

		return c.JSON(http.StatusOK, Response{usr})
	}
}

// @Summary Get Profile Route
// @Description Signing user in
// @ID get-profile
// @Accept json
// @Produce json
// @Success 200 object models.User
// @Failure 401 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/profile [get]
func (uh *UserHandler) GetProfile(c echo.Context) error {
	sess, ok := c.Get("session").(*models.Session)

	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	usr, err := uh.UUsecase.GetByID(sess.UserID)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{err.Error()})
	}

	return c.JSON(http.StatusOK, Response{usr})
}

func (uh *UserHandler) GetAvatar(c echo.Context) error {
	return nil
}

// @Summary Set Avatar Route
// @Description Signing user in
// @ID set-avatar
// @Accept multipart/form-data
// @Produce json
// @Param Data body string true "multipart/form-data"
// @Success 200 object models.User
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/avatar [post]
func (uh *UserHandler) SetAvatar(c echo.Context) error {
	file, err := c.FormFile("avatar")

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{vars.ErrRetrievingError.Error()})
	}

	src, err := file.Open()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	defer src.Close()

	rootPath, _ := os.Getwd()
	avatarsPath := "/resources/avatars/"
	destPath := filepath.Join(rootPath, avatarsPath)

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	bytes, err := ioutil.ReadAll(src)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	fileType := http.DetectContentType(bytes)
	extensions, err := mime.ExtensionsByType(fileType)

	sess, ok := c.Get("session").(*models.Session)

	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	avatarName := filepath.Join(fmt.Sprint(sess.UserID) + "_avatar" + extensions[0])
	destFile, err := os.Create(filepath.Join(destPath, avatarName))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	defer destFile.Close()

	_, err = destFile.Write(bytes)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	usr, err := uh.UUsecase.GetByID(sess.UserID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	usr.Avatar = filepath.Join(avatarsPath, avatarName)

	return c.JSON(http.StatusOK, Response{usr})
}
