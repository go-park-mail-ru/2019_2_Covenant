package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/playlist"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	"2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	. "2019_2_Covenant/tools/vars"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type UserHandler struct {
	base_handler.BaseHandler
	UUsecase user.Usecase
	SUsecase session.Usecase
	PUsecase playlist.Usecase
}

func NewUserHandler(uUC user.Usecase,
	sUC session.Usecase,
	pUC playlist.Usecase,
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
		PUsecase: pUC,
	}
}

func (uh *UserHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/users", uh.CreateUser())

	e.GET("/api/v1/users/:nickname", uh.GetOtherProfile(), uh.MManager.CheckAuthStrictly)
	e.GET("/api/v1/users/:id/subscriptions", uh.GetUserSubscriptions(), uh.MManager.CheckAuthStrictly)
	e.GET("/api/v1/users/:id/playlists", uh.GetUserPlaylists(), uh.MManager.CheckAuthStrictly)

	e.GET("/api/v1/profile", uh.GetProfile(), uh.MManager.CheckAuthStrictly)
	e.PUT("/api/v1/profile", uh.UpdateUser(), uh.MManager.CheckAuthStrictly)
	e.PUT("/api/v1/profile/password", uh.UpdatePassword(), uh.MManager.CheckAuthStrictly)
	e.PUT("/api/v1/profile/avatar", uh.UploadAvatar(), uh.MManager.CheckAuthStrictly)
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

		newUser := models.NewUser(request.Email, request.Nickname, request.Password)

		if err := uh.UUsecase.Store(newUser); err != nil {
			uh.Logger.Log(c, "info", "User store error.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		sess, cookie := models.NewSession(newUser.ID)
		c.SetCookie(cookie)

		if err := uh.SUsecase.Store(sess); err != nil {
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
func (uh *UserHandler) UploadAvatar() echo.HandlerFunc {
	rootPath, _ := os.Getwd()

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			uh.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		file, err := c.FormFile("file")
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

		filePath := strings.ReplaceAll(
			fmt.Sprintf("%s%d-%s", AVATARS_PATH, sess.UserID, file.Filename),
			" ",
			"",
		)

		dest, err := os.Create(filepath.Join(rootPath, filePath))
		if err != nil {
			uh.Logger.Log(c, "error", "Can't create file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		defer dest.Close()

		if _, err = io.Copy(dest, src); err != nil {
			uh.Logger.Log(c, "error", "Can't copy file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		usr, err := uh.UUsecase.UpdateAvatar(sess.UserID, filePath)

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

func (uh *UserHandler) GetOtherProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			uh.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		uNickname := c.Param("nickname")

		usr, err := uh.UUsecase.GetByNickname(uNickname, sess.UserID)

		if err != nil {
			uh.Logger.Log(c, "info", "Error while getting other profile.", err)
			return c.JSON(http.StatusBadRequest, Response{
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

func (uh *UserHandler) GetUserSubscriptions() echo.HandlerFunc {
	type Request struct {
		Count  uint64 `query:"count" validate:"required"`
		Offset uint64 `query:"offset"`
	}

	return func(c echo.Context) error {
		uID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			uh.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := uh.ReqReader.Read(c, request, nil); err != nil {
			uh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		followers, totalFollowers, err := uh.UUsecase.GetFollowers(uint64(uID), request.Count, request.Offset)

		if err != nil {
			uh.Logger.Log(c, "error", "Error while getting followers.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		following, totalFollowing, err := uh.UUsecase.GetFollowing(uint64(uID), request.Count, request.Offset)

		if err != nil {
			uh.Logger.Log(c, "error", "Error while getting following.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"followers": followers,
				"total_followers": totalFollowers,
				"following": following,
				"total_following": totalFollowing,
			},
		})
	}
}

func (uh *UserHandler) GetUserPlaylists() echo.HandlerFunc {
	type Request struct {
		Count  uint64 `query:"count" validate:"required"`
		Offset uint64 `query:"offset"`
	}

	return func(c echo.Context) error {
		uID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			uh.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := uh.ReqReader.Read(c, request, nil); err != nil {
			uh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		playlists, total, err := uh.PUsecase.Fetch(uint64(uID), request.Count, request.Offset)

		if err != nil {
			uh.Logger.Log(c, "error", "Error while getting playlists.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"playlists": playlists,
				"total": total,
			},
		})
	}
}