package middleware

import (
	"2019_2_Covenant/pkg/session"
	"2019_2_Covenant/pkg/user"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type Middleware struct {
	sUC session.Usecase
	uUC user.Usecase
}

func NewMiddleware(uUsecase user.Usecase, sUsecase session.Usecase) Middleware {
	return Middleware{
		sUC: sUsecase,
		uUC: uUsecase,
	}
}

func (m Middleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
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
