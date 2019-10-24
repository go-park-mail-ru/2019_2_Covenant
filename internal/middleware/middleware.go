package middleware

import (
	"2019_2_Covenant/internal/session"
	user2 "2019_2_Covenant/internal/user"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MiddlewareManager struct {
	sUC session.Usecase
	uUC user2.Usecase
}

func NewMiddlewareManager(uUsecase user2.Usecase, sUsecase session.Usecase) MiddlewareManager {
	return MiddlewareManager{
		sUC: sUsecase,
		uUC: uUsecase,
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
