package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	"2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	. "2019_2_Covenant/tools/vars"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UserHandler struct {
	base_handler.BaseHandler
	UUsecase user.Usecase
	SUsecase session.Usecase
}

func NewUserHandler(uUC user.Usecase,
	sUC session.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *UserHandler {
	return &UserHandler{
		BaseHandler: base_handler.BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		UUsecase: uUC,
		SUsecase: sUC,
	}
}

func (uh *UserHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/users", uh.CreateUser())

	e.GET("/api/v1/profile", uh.GetProfile(), uh.MManager.CheckAuth)
	e.POST("/api/v1/profile", uh.UpdateUser(), uh.MManager.CheckAuth)
	e.POST("/api/v1/profile/password", uh.UpdatePassword(), uh.MManager.CheckAuth)
	e.POST("/api/v1/profile/avatar", uh.SetAvatar(), uh.MManager.CheckAuth)
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
// @Router /api/v1/users [post]
func (uh *UserHandler) CreateUser() echo.HandlerFunc {
	type Request struct {
		Nickname         string `json:"nickname" validate:"required"`
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,gte=6"`
		PassConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	}

	correctData := func(req interface{}) bool {
		return strings.Contains(req.(*Request).Password, " ") == false &&
			strings.Contains(req.(*Request).Nickname, " ") == false
	}

	return func(c echo.Context) error {
		request := &Request{}

		if err := uh.ReqReader.Read(c, request, correctData); err != nil {
			uh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		usr, err := uh.UUsecase.GetByEmail(request.Email)

		if err == nil {
			uh.Logger.Log(c, "info", "Already exists.", "User ID:", usr.ID)
			return c.JSON(http.StatusBadRequest, Response{
				Error: ErrAlreadyExist.Error(),
			})
		}

		newUser := models.NewUser(request.Email, request.Nickname, request.Password)

		if err = uh.UUsecase.Store(newUser); err != nil {
			uh.Logger.Log(c, "error", "User store error.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		sess, cookie := models.NewSession(newUser.ID)
		c.SetCookie(cookie)

		if err = uh.SUsecase.Store(sess); err != nil {
			uh.Logger.Log(c, "error", "Session store error.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		token, err := models.NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, time.Now().Add(24*time.Hour))
		c.Response().Header().Set("X-CSRF-Token", token)

		if err != nil {
			uh.Logger.Log(c, "error", "CSRF Token generating error.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"user": newUser,
			},
		})
	}
}

// @Tags Profile
// @Summary Edit Profile Route
// @Description Edit user profile
// @ID edit-profile
// @Accept json
// @Produce json
// @Param Data body object true "JSON that contains user data to edit"
// @Success 200 object models.User
// @Failure 400 object Response
// @Failure 401 object Response
// @Failure 409 object Response
// @Failure 500 object Response
// @Router /api/v1/profile [post]
func (uh *UserHandler) UpdateUser() echo.HandlerFunc {
	type Request struct {
		Email    string `json:"email" validate:"required,email"`
		Nickname string `json:"nickname" validate:"required"`
	}

	correctData := func(req interface{}) bool {
		return strings.Contains(req.(*Request).Nickname, " ") == false
	}

	return func(c echo.Context) error {
		usr, ok := c.Get("user").(*models.User)

		if !ok {
			uh.Logger.Log(c, "error", "Can't extract user from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := uh.ReqReader.Read(c, request, correctData); err != nil {
			uh.Logger.Log(c, "info", "Invalid request:", request)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		usr, err := uh.UUsecase.Update(usr.ID, request.Nickname, request.Email)

		if err != nil {
			uh.Logger.Log(c, "info", "Error while updating user data.", err)
			return c.JSON(http.StatusConflict, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"user": usr,
			},
		})
	}
}

// @Tags Profile
// @Summary Edit Profile Route
// @Description Edit user profile
// @ID edit-profile
// @Accept json
// @Produce json
// @Param Data body object true "JSON that contains user data to edit"
// @Success 200 object models.User
// @Failure 400 object Response
// @Failure 401 object Response
// @Failure 500 object Response
// @Router /api/v1/profile/password [post]
func (uh *UserHandler) UpdatePassword() echo.HandlerFunc {
	type Request struct {
		OldPassword      string `json:"old_password" validate:"required"`
		Password         string `json:"password" validate:"required"`
		PassConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	}

	correctData := func(req interface{}) bool {
		return strings.Contains(req.(*Request).Password, " ") == false
	}

	return func(c echo.Context) error {
		usr, ok := c.Get("user").(*models.User)

		if !ok {
			uh.Logger.Log(c, "error", "Can't extract user from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := uh.ReqReader.Read(c, request, correctData); err != nil {
			uh.Logger.Log(c, "info", "Invalid request:", request)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if !usr.Verify(request.OldPassword) {
			uh.Logger.Log(c, "info", "Bad old password.", "User:", usr.Nickname)
			return c.JSON(http.StatusBadRequest, Response{
				Error: ErrBadParam.Error(),
			})
		}

		if err := uh.UUsecase.UpdatePassword(usr.ID, request.Password); err != nil {
			uh.Logger.Log(c, "info", "Error while updating user data.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

// @Tags Profile
// @Summary Get Profile Route
// @Description Get user profile
// @ID get-profile
// @Accept json
// @Produce json
// @Success 200 object models.User
// @Failure 401 object Response
// @Failure 500 object Response
// @Router /api/v1/profile [get]
func (uh *UserHandler) GetProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		usr, ok := c.Get("user").(*models.User)

		if !ok {
			uh.Logger.Log(c, "error", "Can't extract user from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"user": usr,
			},
		})
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
// @Failure 400 object Response
// @Failure 401 object Response
// @Failure 500 object Response
// @Router /api/v1/profile/avatar [post]
func (uh *UserHandler) SetAvatar() echo.HandlerFunc {
	rootPath, _ := os.Getwd()
	avatarsPath := "/resources/avatars/"
	destPath := filepath.Join(rootPath, avatarsPath)

	return func(c echo.Context) error {
		file, err := c.FormFile("avatar")

		if err != nil {
			uh.Logger.Log(c, "info", "Can't extract file from request.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: ErrRetrievingError.Error(),
			})
		}

		src, err := file.Open()

		if err != nil {
			uh.Logger.Log(c, "error", "Can't open file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		defer src.Close()

		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			uh.Logger.Log(c, "error", "There is no dir for avatars.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		bytes, err := ioutil.ReadAll(src)

		if err != nil {
			uh.Logger.Log(c, "error", "Can't read file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		fileType := http.DetectContentType(bytes)
		extensions, _ := mime.ExtensionsByType(fileType)

		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			uh.Logger.Log(c, "info", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		avatarName := filepath.Join(fmt.Sprint(sess.UserID) + "_avatar" + extensions[0])
		destFile, err := os.Create(filepath.Join(destPath, avatarName))

		if err != nil {
			uh.Logger.Log(c, "error", "Can't create avatar file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		defer destFile.Close()

		_, err = destFile.Write(bytes)

		if err != nil {
			uh.Logger.Log(c, "error", "Error while writing bytes in destFile.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		usr, err := uh.UUsecase.UpdateAvatar(sess.UserID, filepath.Join(avatarsPath, avatarName))

		if err != nil {
			uh.Logger.Log(c, "error", "Error while updating user avatar.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"user": usr,
			},
		})
	}
}
