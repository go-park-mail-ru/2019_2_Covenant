package middleware

import (
	"2019_2_Covenant/pkg/session"
	"2019_2_Covenant/pkg/user"
	"fmt"
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
			err = next(c)
			return err
		}

		sess, err := m.sUC.Get(cookie.Value)

		if err != nil {
			err = next(c)
			return err
		}

		usr, err := m.uUC.GetByID(sess.UserID)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		fmt.Println("authorized")
		return c.JSON(http.StatusOK, usr)
	}
}
