package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/vars"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type SessionHandler struct {
	SUsecase session.Usecase
	MManager middlewares.MiddlewareManager
	Logger   *logrus.Logger
}

func NewSessionHandler(sUC session.Usecase, mManager middlewares.MiddlewareManager, logger *logrus.Logger) *SessionHandler {
	return &SessionHandler{
		SUsecase: sUC,
		MManager: mManager,
		Logger:   logger,
	}
}

func (sh *SessionHandler) Configure(e *echo.Echo) {
	e.GET("/api/v1/get_csrf", sh.GetCSRF(), sh.MManager.CheckAuth)
}

func (sh *SessionHandler) log(c echo.Context, logType string, msg ...interface{}) {
	fields := logrus.Fields{
		"Request Method": c.Request().Method,
		"Remote Address": c.Request().RemoteAddr,
		"Message":        msg,
	}

	switch logType {
	case "error":
		sh.Logger.WithFields(fields).Error(c.Request().URL.Path)
	case "info":
		sh.Logger.WithFields(fields).Info(c.Request().URL.Path)
	case "warning":
		sh.Logger.WithFields(fields).Warning(c.Request().URL.Path)
	}
}

func (sh *SessionHandler) GetCSRF() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			sh.log(c, "info", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		token, err := models.NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, time.Now().Add(24*time.Hour))
		c.Response().Header().Set("X-CSRF-Token", token)

		if err != nil {
			sh.log(c, "error", "CSRF Token generating error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "success",
		})
	}
}
