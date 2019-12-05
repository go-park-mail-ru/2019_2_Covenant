package delivery

import (
	"2019_2_Covenant/internal/follower"
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	. "2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	. "2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type FollowerHandler struct {
	BaseHandler
	FUsecase follower.Usecase
}

func NewFollowerHandler(fUC follower.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *FollowerHandler {
	return &FollowerHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		FUsecase: fUC,
	}
}

func (fh *FollowerHandler) Configure(e *echo.Echo) {
	e.GET("/api/v1/profile/:id", fh.GetProfile(), fh.MManager.CheckAuth)
	e.POST("/api/v1/profile/:id", fh.Subscribe(), fh.MManager.CheckAuth)
	e.DELETE("/api/v1/profile/:id", fh.Unsubscribe(), fh.MManager.CheckAuth)
}

func (fh *FollowerHandler) Subscribe() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			fh.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		uID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			fh.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		if err := fh.FUsecase.Subscribe(uint64(uID), sess.UserID); err != nil {
			fh.Logger.Log(c, "error", "Error while subscribing.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrAlreadyExist.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (fh *FollowerHandler) Unsubscribe() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			fh.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		uID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			fh.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		if err := fh.FUsecase.Unsubscribe(uint64(uID), sess.UserID); err != nil {
			fh.Logger.Log(c, "info", "Error while unsubscribing.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (fh *FollowerHandler) GetProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		uID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			fh.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		usr, err := fh.FUsecase.GetProfile(uint64(uID))

		if err != nil {
			fh.Logger.Log(c, "info", "Error while unsubscribing.", err)
			return c.JSON(http.StatusBadRequest, Response{
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
