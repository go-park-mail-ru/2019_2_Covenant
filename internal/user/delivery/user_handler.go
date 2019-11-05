package delivery

import (
	_middleware "2019_2_Covenant/internal/middleware"
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
	MManager _middleware.MiddlewareManager
}

func NewUserHandler(uUC user.Usecase, sUC session.Usecase, mManager _middleware.MiddlewareManager) *UserHandler {
	return &UserHandler{
		UUsecase: uUC,
		SUsecase: sUC,
		MManager: mManager,
	}
}

func (uh *UserHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/signup", uh.SignUp())
	e.POST("/api/v1/login", uh.LogIn())
	e.POST("/api/v1/profile", uh.EditProfile(), uh.MManager.CheckAuth)
	e.GET("/api/v1/profile", uh.GetProfile(), uh.MManager.CheckAuth)
	e.POST("/api/v1/avatar", uh.SetAvatar(), uh.MManager.CheckAuth)
	e.GET("/api/v1/avatar", uh.GetAvatar(), uh.MManager.CheckAuth)
	e.GET("/api/v1/logout", uh.LogOut(), uh.MManager.CheckAuth)
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

// @Tags User
// @Summary SignUp Route
// @Description Signing user up
// @ID sign-up-user
// @Accept json
// @Produce json
// @Param Data body object true "JSON that contains user sign up data"
// @Success 200 object models.User
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/signup [post]
func (uh *UserHandler) SignUp() echo.HandlerFunc {
	type UserReg struct {
		Nickname string `json:"nickname" validate:"required"`
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
			Email:         userRegData.Email,
			PlainPassword: userRegData.Password,
			Nickname:      userRegData.Nickname,
		}

		usr, err = uh.UUsecase.Store(newUser)

		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
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

		return c.JSON(http.StatusOK, Response{newUser})
	}
}

// @Tags User
// @Summary LogIn Route
// @Description Logging user in
// @ID log-in-user
// @Accept json
// @Produce json
// @Param Data body object true "JSON that contains user login data"
// @Success 200 object models.User
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/login [post]
func (uh *UserHandler) LogIn() echo.HandlerFunc {
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

		if !usr.Verify(userLoginData.Password) {
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

// @Tags User
// @Summary Edit Profile Route
// @Description Edit user profile
// @ID edit-profile
// @Accept json
// @Produce json
// @Param Data body object true "JSON that contains user data to edit"
// @Success 200 object models.User
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/profile [post]
func (uh *UserHandler) EditProfile() echo.HandlerFunc {
	type UserEdit struct {
		Nickname string `json:"nickname" validate:"required"`
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

		if usr, err = uh.UUsecase.UpdateNickname(usr.ID, userEditData.Nickname); err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseError{err.Error()})
		}

		return c.JSON(http.StatusOK, Response{usr})
	}
}

// @Tags User
// @Summary Get Profile Route
// @Description Get user profile
// @ID get-profile
// @Accept json
// @Produce json
// @Success 200 object models.User
// @Failure 401 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/profile [get]
func (uh *UserHandler) GetProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
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
}

// @Tags User
// @Summary Get Avatar Route
// @Description Signing user in
// @ID get-avatar
// @Accept json
// @Produce json
// @Param Data body string true "multipart/form-data"
// @Success 200 object models.User
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/avatar [get]
func (uh *UserHandler) GetAvatar() echo.HandlerFunc {
	rootPath, _ := os.Getwd()

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		usr, err := uh.UUsecase.GetByID(sess.UserID)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, ResponseError{err.Error()})
		}

		avatarPath := usr.Avatar

		destPath := filepath.Join(rootPath, avatarPath)
		//src, err := os.Open(destPath)

		//defer src.Close()

		return c.File(destPath)
	}
}

// @Tags User
// @Summary Set Avatar Route
// @Description Set user avatar
// @ID set-avatar
// @Accept multipart/form-data
// @Produce json
// @Param Data body string true "multipart/form-data"
// @Success 200 object models.User
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/avatar [post]
func (uh *UserHandler) SetAvatar() echo.HandlerFunc {
	rootPath, _ := os.Getwd()
	avatarsPath := "/resources/avatars/"
	destPath := filepath.Join(rootPath, avatarsPath)

	return func(c echo.Context) error {
		file, err := c.FormFile("avatar")

		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseError{vars.ErrRetrievingError.Error()})
		}

		src, err := file.Open()

		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		defer src.Close()

		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		bytes, err := ioutil.ReadAll(src)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		fileType := http.DetectContentType(bytes)
		extensions, _ := mime.ExtensionsByType(fileType)

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

		usr, err := uh.UUsecase.UpdateAvatar(sess.UserID, filepath.Join(avatarsPath, avatarName))

		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		return c.JSON(http.StatusOK, Response{usr})
	}
}

// @Tags User
// @Summary Log Out Route
// @Description Logging user out
// @ID log-out-user
// @Accept json
// @Produce json
// @Success 200 object models.User
// @Failure 400 object ResponseError
// @Failure 404 object ResponseError
// @Failure 500 object ResponseError
// @Router /api/v1/logout [get]
func (uh *UserHandler) LogOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		if err := uh.SUsecase.DeleteByID(sess.ID); err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
		}

		cookie := &http.Cookie{
			Name:    "Covenant",
			Value:   sess.Data,
			Expires: time.Now().AddDate(0, 0, -1),
		}

		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]string{
			"message": "logout",
		})
	}
}
