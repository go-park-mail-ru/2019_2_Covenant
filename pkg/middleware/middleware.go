package middleware

import (
	"2019_2_Covenant/pkg/session"
	"2019_2_Covenant/pkg/user"
	"github.com/labstack/echo"
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

func (m MiddlewareManager) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
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
