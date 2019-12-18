package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	. "2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	"2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type SessionHandler struct {
	BaseHandler
	SUsecase  session.Usecase
	UUsecase  user.Usecase
}

func NewSessionHandler(sUC session.Usecase,
	uUC user.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *SessionHandler {
	return &SessionHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		SUsecase: sUC,
		UUsecase: uUC,
	}
}

func (sh *SessionHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/session", sh.CreateSession())
	e.DELETE("/api/v1/session", sh.DeleteSession(), sh.MManager.CheckAuthStrictly)

	e.GET("/api/v1/csrf", sh.GetCSRF(), sh.MManager.CheckAuthStrictly)
}

// @Tags Session
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
// @Router /api/v1/session [post]
func (sh *SessionHandler) CreateSession() echo.HandlerFunc {
	type Request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	return func(c echo.Context) error {
		request := &Request{}

		if err := sh.ReqReader.Read(c, request, nil); err != nil {
			sh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		usr, err := sh.UUsecase.GetByEmail(request.Email)

		if err != nil {
			sh.Logger.Log(c, "info", "Error while getting user by EMAIL.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if !usr.Verify(request.Password) {
			sh.Logger.Log(c, "info", "Bad authentication.", "User:", usr.Nickname)
			return c.JSON(http.StatusBadRequest, Response{
				Error: vars.ErrBadParam.Error(),
			})
		}

		sess, cookie := models.NewSession(usr.ID)
		c.SetCookie(cookie)

		if err = sh.SUsecase.Store(sess); err != nil {
			sh.Logger.Log(c, "error", "Session store error.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		token, err := models.NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, time.Now().Add(24*time.Hour))
		c.Response().Header().Set("X-CSRF-Token", token)

		if err != nil {
			sh.Logger.Log(c, "error", "CSRF Token generating error.", err)
			return c.JSON(http.StatusInternalServerError, Response{
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

// @Tags Session
// @Summary Log Out Route
// @Description Logging user out
// @ID log-out-user
// @Accept json
// @Produce json
// @Success 200 object models.User
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
// @Router /api/v1/session [delete]
func (sh *SessionHandler) DeleteSession() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			sh.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		if err := sh.SUsecase.DeleteByID(sess.ID); err != nil {
			sh.Logger.Log(c, "error", "Error while deleting session.", err.Error())
			return c.JSON(http.StatusNotFound, Response{
				Error: err.Error(),
			})
		}

		cookie := &http.Cookie{
			Name:    "Covenant",
			Value:   sess.Data,
			Expires: time.Now().AddDate(0, 0, -1),
		}

		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (sh *SessionHandler) GetCSRF() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			sh.Logger.Log(c, "info", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		token, err := models.NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, time.Now().Add(24*time.Hour))
		c.Response().Header().Set("X-CSRF-Token", token)

		if err != nil {
			sh.Logger.Log(c, "error", "CSRF Token generating error.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}
