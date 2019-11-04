package middleware

import (
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

		_, err = m.uUC.GetByID(sess.UserID)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		c.Set("session", sess)

		return next(c)
	}
}
