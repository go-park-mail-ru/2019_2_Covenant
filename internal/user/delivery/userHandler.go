package delivery

import (
	"2019_2_Covenant/internal/middleware"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	user2 "2019_2_Covenant/internal/user"
	vars2 "2019_2_Covenant/internal/vars"
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
	UUsecase user2.Usecase
	SUsecase session.Usecase
	MManager middleware.MiddlewareManager
}

func NewUserHandler(uUC user2.Usecase, sUC session.Usecase) *UserHandler {
	return &UserHandler{
		UUsecase: uUC,
		SUsecase: sUC,
		MManager: middleware.NewMiddlewareManager(uUC, sUC),
	}
}

func (uh *UserHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/signup", uh.SignUp)
	e.POST("/api/v1/signin", uh.SignIn)
	e.POST("/api/v1/profile", uh.Profile, uh.MManager.CheckAuth)
	e.GET("/api/v1/profile", uh.Profile, uh.MManager.CheckAuth)
	e.POST("/api/v1/avatar", uh.Avatar, uh.MManager.CheckAuth)
	e.GET("/api/v1/avatar", uh.Avatar, uh.MManager.CheckAuth)
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
		return false, vars2.ErrBadParam
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
func (uh *UserHandler) SignUp(c echo.Context) error {
	var userRegData models.UserReg
	err := c.Bind(&userRegData)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{err.Error()})
	}

	if ok, err := isValidRequest(userRegData); !ok {
		return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
	}

	usr, err := uh.UUsecase.GetByEmail(userRegData.Email)

	if usr != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{vars2.ErrAlreadyExist.Error()})
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

	err = uh.SUsecase.Store(sess)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, Response{newUser})
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
func (uh *UserHandler) SignIn(c echo.Context) error {
	var userLoginData models.UserLogin
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
		return c.JSON(http.StatusBadRequest, ResponseError{vars2.ErrBadParam.Error()})
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
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, Response{usr})
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
func (uh *UserHandler) editProfile(c echo.Context) error {
	sess, ok := c.Get("session").(*models.Session)

	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	usr, err := uh.UUsecase.GetByID(sess.UserID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
	}

	var userEditData models.UserEdit
	err = c.Bind(&userEditData)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{err.Error()})
	}

	if ok, err := isValidRequest(userEditData); !ok {
		return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
	}

	usr.Name = userEditData.Name
	usr.Surname = userEditData.Surname

	return c.JSON(http.StatusOK, Response{usr})
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
func (uh *UserHandler) getProfile(c echo.Context) error {
	sess, ok := c.Get("session").(*models.Session)

	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	usr, err := uh.UUsecase.GetByID(sess.UserID)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseError{err.Error()})
	}

	return c.JSON(http.StatusOK, Response{usr})
}

func (uh *UserHandler) Profile(c echo.Context) error {
	var err error

	switch c.Request().Method {
	case echo.GET:
		err = uh.getProfile(c)
	case echo.POST:
		err = uh.editProfile(c)
	default:
		err = nil
	}

	return err
}

func (uh *UserHandler) getAvatar(c echo.Context) error {
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
func (uh *UserHandler) setAvatar(c echo.Context) error {
	file, err := c.FormFile("avatar")

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{vars2.ErrRetrievingError.Error()})
	}

	src, err := file.Open()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	defer src.Close()

	rootPath, _ := os.Getwd()
	avatarsPath := "/resources/avatars/"
	destPath := filepath.Join(rootPath, avatarsPath)

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	bytes, err := ioutil.ReadAll(src)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	fileType := http.DetectContentType(bytes)
	extensions, err := mime.ExtensionsByType(fileType)

	sess, ok := c.Get("session").(*models.Session)

	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	avatarName := filepath.Join(fmt.Sprint(sess.UserID) + "_avatar" + extensions[0])
	destFile, err := os.Create(filepath.Join(destPath, avatarName))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	defer destFile.Close()

	_, err = destFile.Write(bytes)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	usr, err := uh.UUsecase.GetByID(sess.UserID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars2.ErrInternalServerError.Error()})
	}

	usr.Avatar = filepath.Join(avatarsPath, avatarName)

	return c.JSON(http.StatusOK, Response{usr})
}

func (uh *UserHandler) Avatar(c echo.Context) error {
	var err error

	switch c.Request().Method {
	case echo.GET:
		err = uh.getAvatar(c)
	case echo.POST:
		err = uh.setAvatar(c)
	default:
		err = nil
	}

	return err
}
