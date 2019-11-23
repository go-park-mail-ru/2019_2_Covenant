package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/vars"
	"2019_2_Covenant/pkg/logger"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	SUsecase session.Usecase
	MManager middlewares.MiddlewareManager
	Logger   *logger.LogrusLogger
}

func NewSessionHandler(sUC session.Usecase,
	mManager middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *SessionHandler {
	return &SessionHandler{
		SUsecase: sUC,
		MManager: mManager,
		Logger:   logger,
	}
}

func (sh *SessionHandler) Configure(e *echo.Echo) {
	e.GET("/api/v1/get_csrf", sh.GetCSRF(), sh.MManager.CheckAuth)
}

func (sh *SessionHandler) GetCSRF() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			sh.Logger.Log(c, "info", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		token, err := models.NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, time.Now().Add(24*time.Hour))
		c.Response().Header().Set("X-CSRF-Token", token)

		if err != nil {
			sh.Logger.Log(c, "error", "CSRF Token generating error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "success",
		})
	}
}
