package delivery

import (
	"2019_2_Covenant/pkg/middleware"
	"2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/session"
	"2019_2_Covenant/pkg/user"
	"2019_2_Covenant/pkg/vars"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
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
	//e.Use(uh.MManager.CheckAuth)
	e.POST("/api/v1/signup", uh.SignUp)
	e.POST("/api/v1/signin", uh.SignIn)
	e.POST("/api/v1/profile", uh.Profile, uh.MManager.CheckAuth)
	e.GET("/api/v1/profile", uh.Profile, uh.MManager.CheckAuth)
}

type ResponseError struct {
	Error string `json:"error"`
}

func isValidRequest(usr interface{}) (bool, error) {
	v := validator.New()
	err := v.Struct(usr)

	if err != nil {
		return false, vars.ErrBadParam
	}

	return true, nil
}

// curl -X POST 127.0.0.1:8000/api/v1/signup -H 'Content-Type: application/json' \
// -d '{"email": "m@mail.ru", "username": "Marshal", "password": "12345312"}'

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

	err = uh.SUsecase.Store(sess)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, newUser)
}

// curl -X POST 127.0.0.1:8000/api/v1/signin -H 'Content-Type: application/json' \
// -d '{"email": "m@mail.ru", "password": "12345312"}'

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

	return c.JSON(http.StatusOK, usr)
}

func (uh *UserHandler) editProfile(c echo.Context) error {
	sess, ok := c.Get("session").(*models.Session)

	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
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

	return c.JSON(http.StatusOK, usr)
}

func (uh *UserHandler) getProfile(c echo.Context) error {
	sess, ok := c.Get("session").(models.Session)

	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{vars.ErrInternalServerError.Error()})
	}

	usr, err := uh.UUsecase.GetByID(sess.UserID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{err.Error()})
	}

	return c.JSON(http.StatusOK, usr)
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

func (uh *UserHandler) setAvatar(c echo.Context) error {
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

	return c.JSON(http.StatusOK, usr)
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
