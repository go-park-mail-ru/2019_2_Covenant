package middlewares

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/internal/vars"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MiddlewareManager struct {
	sUC session.Usecase
	uUC user.Usecase
}

func NewMiddlewareManager(uUsecase user.Usecase, sUsecase session.Usecase) MiddlewareManager {
	return MiddlewareManager{
		sUC: sUsecase,
		uUC: uUsecase,
	}
}

func (m *MiddlewareManager) CSRFCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("X-Csrf-Token")

		usr := c.Get("user").(*models.User)
		sess := c.Get("session").(*models.Session)

		ok, err := models.NewCSRFTokenManager("Covenant").Verify(usr, sess, token)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		if !ok {
			return c.JSON(http.StatusBadRequest, vars.ResponseError{
				Error: vars.ErrExpired.Error(),
			})
		}

		return next(c)
	}
}

func (m *MiddlewareManager) CORSMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "*")
		c.Response().Header().Set("Access-Control-Allow-Origin", "http://front.covenant.fun:8000/")
		c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

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
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "unauthorized",
			})
		}

		sess, err := m.sUC.Get(cookie.Value)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "unauthorized",
			})
		}

		usr, err := m.uUC.GetByID(sess.UserID)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		c.Set("session", sess)
		c.Set("user", usr)

		return next(c)
	}
}
