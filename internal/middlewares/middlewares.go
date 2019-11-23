package middlewares

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/internal/vars"
	"2019_2_Covenant/pkg/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type MiddlewareManager struct {
	sUC    session.Usecase
	uUC    user.Usecase
	logger *logger.LogrusLogger
}

func NewMiddlewareManager(uUsecase user.Usecase,
	sUsecase session.Usecase,
	logger *logger.LogrusLogger) MiddlewareManager {
	return MiddlewareManager{
		sUC:    sUsecase,
		uUC:    uUsecase,
		logger: logger,
	}
}

func (m *MiddlewareManager) CSRFCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("X-Csrf-Token")

		sess := c.Get("session").(*models.Session)

		ok, err := models.NewCSRFTokenManager("Covenant").Verify(sess.UserID, sess.Data, token)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		if !ok {
			return c.JSON(http.StatusBadRequest, vars.ResponseError{
				Error: vars.ErrBadCSRF.Error(),
			})
		}

		return next(c)
	}
}

func (m *MiddlewareManager) AccessLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		m.logger.L.WithFields(logrus.Fields{
			"Request Method": c.Request().Method,
			"Remote Address": c.Request().RemoteAddr,
			"Work Time":      time.Since(start),
		}).Info(c.Request().URL.Path)

		return next(c)
	}
}

func (m *MiddlewareManager) CORSMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		origin := c.Request().Header.Get("Origin")

		if origin == "http://localhost:3000" || origin == "http://front.covenant.fun:3000" || origin == "http://front.covenant.fun:5000"{
			c.Response().Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request().Method == "OPTIONS" {
			return nil
		}

		return next(c)
	}
}

func (m *MiddlewareManager) PanicRecovering(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				c.Error(vars.ErrInternalServerError)
			}
		}()

		return next(c)
	}
}

func (m *MiddlewareManager) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Covenant")

		if err != nil {
			return c.JSON(http.StatusUnauthorized, vars.ResponseError{
				Error:vars.ErrUnathorized.Error(),
			})
		}

		sess, err := m.sUC.Get(cookie.Value)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, vars.ResponseError{
				Error:vars.ErrUnathorized.Error(),
			})
		}

		c.Set("session", sess)

		return next(c)
	}
}
