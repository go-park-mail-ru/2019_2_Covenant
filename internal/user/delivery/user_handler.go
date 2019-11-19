package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/internal/vars"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

// Для тестирования только этого файла:
// go test -v -cover -race ./internal/user/delivery

type UserHandler struct {
	UUsecase user.Usecase
	SUsecase session.Usecase
	MManager middlewares.MiddlewareManager
	Logger   *logrus.Logger
}

func NewUserHandler(uUC user.Usecase, sUC session.Usecase, mManager middlewares.MiddlewareManager, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		UUsecase: uUC,
		SUsecase: sUC,
		MManager: mManager,
		Logger:   logger,
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

func isValidRequest(usr interface{}) error {
	v := validator.New()
	err := v.Struct(usr)

	if err != nil {
		return vars.ErrBadParam
	}

	return nil
}

func (uh *UserHandler) log(c echo.Context, logType string, msg ...interface{}) {
	fields := logrus.Fields{
		"Request Method": c.Request().Method,
		"Remote Address": c.Request().RemoteAddr,
		"Message":        msg,
	}

	switch logType {
	case "error":
		uh.Logger.WithFields(fields).Error(c.Request().URL.Path)
	case "info":
		uh.Logger.WithFields(fields).Info(c.Request().URL.Path)
	case "warning":
		uh.Logger.WithFields(fields).Warning(c.Request().URL.Path)
	}
}

// @Tags User
// @Summary SignUp Route
// @Description Signing user up
// @ID sign-up-user
// @Accept json
// @Produce json
// @Param Data body object true "JSON that contains user sign up data"
// @Success 200 object models.User
// @Failure 400 object vars.ResponseError
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
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
			uh.log(c, "error", "Can't read request body.")
			return c.JSON(http.StatusUnprocessableEntity, vars.ResponseError{Error: err.Error()})
		}

		if err := isValidRequest(userRegData); err != nil {
			uh.log(c, "info", "Invalid request.", userRegData)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
		}

		usr, err := uh.UUsecase.GetByEmail(userRegData.Email)

		if usr != nil {
			uh.log(c, "info", "Already exist.", "User ID:", usr.ID)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{
				Error: vars.ErrAlreadyExist.Error(),
			})
		}

		newUser := &models.User{
			Email:         userRegData.Email,
			PlainPassword: userRegData.Password,
			Nickname:      userRegData.Nickname,
		}

		usr, err = uh.UUsecase.Store(newUser)

		if err != nil {
			uh.log(c, "error", "User store error.", err)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
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
			uh.log(c, "error", "Session store error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		token, err := models.NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, time.Now().Add(24*time.Hour))
		c.Response().Header().Set("X-CSRF-Token", token)

		if err != nil {
			uh.log(c, "error", "CSRF Token generating error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, vars.Response{Body: newUser})
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
// @Failure 400 object vars.ResponseError
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
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
			uh.log(c, "error", "Can't read request body.")
			return c.JSON(http.StatusUnprocessableEntity, vars.ResponseError{Error: err.Error()})
		}

		if err := isValidRequest(userLoginData); err != nil {
			uh.log(c, "info", "Invalid request.", userLoginData)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
		}

		usr, err := uh.UUsecase.GetByEmail(userLoginData.Email)

		if usr == nil {
			uh.log(c, "info", "Error while getting user by EMAIL.", err)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
		}

		if !usr.Verify(userLoginData.Password) {
			uh.log(c, "info", "Bad authentication.", "User ID:", usr.Nickname)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{
				Error: vars.ErrBadParam.Error(),
			})
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
			uh.log(c, "error", "Session store error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		token, err := models.NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, time.Now().Add(24*time.Hour))
		c.Response().Header().Set("X-CSRF-Token", token)

		if err != nil {
			uh.log(c, "error", "CSRF Token generating error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, vars.Response{Body: usr})
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
// @Failure 400 object vars.ResponseError
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
// @Router /api/v1/profile [post]
func (uh *UserHandler) EditProfile() echo.HandlerFunc {
	type UserEdit struct {
		Nickname string `json:"nickname" validate:"required"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)
		if !ok {
			uh.log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		usr, err := uh.UUsecase.GetByID(sess.UserID)

		if err != nil {
			uh.log(c, "info", "Error while getting user by ID.", err)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
		}

		var userEditData UserEdit
		err = c.Bind(&userEditData)

		if err != nil {
			uh.log(c, "error", "Can't read request body.")
			return c.JSON(http.StatusUnprocessableEntity, vars.ResponseError{Error: err.Error()})
		}

		if err := isValidRequest(userEditData); err != nil {
			uh.log(c, "info", "Invalid request.", userEditData)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
		}

		if usr, err = uh.UUsecase.UpdateNickname(usr.ID, userEditData.Nickname); err != nil {
			uh.log(c, "error", "Error while updating user nickname.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, vars.Response{Body: usr})
	}
}

// @Tags User
// @Summary Get Profile Route
// @Description Get user profile
// @ID get-profile
// @Accept json
// @Produce json
// @Success 200 object models.User
// @Failure 401 object vars.ResponseError
// @Failure 500 object vars.ResponseError
// @Router /api/v1/profile [get]
func (uh *UserHandler) GetProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			uh.log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		usr, err := uh.UUsecase.GetByID(sess.UserID)

		if err != nil {
			uh.log(c, "info", "Error while getting user by ID.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, vars.Response{Body: usr})
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
// @Failure 400 object vars.ResponseError
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
// @Router /api/v1/avatar [get]
func (uh *UserHandler) GetAvatar() echo.HandlerFunc {
	rootPath, _ := os.Getwd()

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			uh.log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		usr, err := uh.UUsecase.GetByID(sess.UserID)

		if err != nil {
			uh.log(c, "info", "Error while getting user by ID.", err)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
		}

		avatarPath := usr.Avatar

		destPath := filepath.Join(rootPath, avatarPath)

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
// @Failure 400 object vars.ResponseError
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
// @Router /api/v1/avatar [post]
func (uh *UserHandler) SetAvatar() echo.HandlerFunc {
	rootPath, _ := os.Getwd()
	avatarsPath := "/resources/avatars/"
	destPath := filepath.Join(rootPath, avatarsPath)

	return func(c echo.Context) error {
		file, err := c.FormFile("avatar")

		if err != nil {
			uh.log(c, "info", "Can't extract file from request.", err)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{
				Error: vars.ErrRetrievingError.Error(),
			})
		}

		src, err := file.Open()

		if err != nil {
			uh.log(c, "error", "Can't open file.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		defer src.Close()

		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			uh.log(c, "error", "There is no dir for avatars.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		bytes, err := ioutil.ReadAll(src)

		if err != nil {
			uh.log(c, "error", "Can't read file.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		fileType := http.DetectContentType(bytes)
		extensions, _ := mime.ExtensionsByType(fileType)

		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			uh.log(c, "info", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		avatarName := filepath.Join(fmt.Sprint(sess.UserID) + "_avatar" + extensions[0])
		destFile, err := os.Create(filepath.Join(destPath, avatarName))

		if err != nil {
			uh.log(c, "error", "Can't create avatar file.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		defer destFile.Close()

		_, err = destFile.Write(bytes)

		if err != nil {
			uh.log(c, "error", "Error while writing bytes in destFile.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		usr, err := uh.UUsecase.UpdateAvatar(sess.UserID, filepath.Join(avatarsPath, avatarName))

		if err != nil {
			uh.log(c, "error", "Error while updating user avatar.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		return c.JSON(http.StatusOK, vars.Response{Body: usr})
	}
}

// @Tags User
// @Summary Log Out Route
// @Description Logging user out
// @ID log-out-user
// @Accept json
// @Produce json
// @Success 200 object models.User
// @Failure 400 object vars.ResponseError
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
// @Router /api/v1/logout [get]
func (uh *UserHandler) LogOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			uh.log(c, "info", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		if err := uh.SUsecase.DeleteByID(sess.ID); err != nil {
			uh.log(c, "error", "Error while deleting session.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
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
