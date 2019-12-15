package delivery

import (
	"2019_2_Covenant/internal/subscriptions"
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	. "2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	. "2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SubscriptionHandler struct {
	BaseHandler
	SUsecase subscriptions.Usecase
	UUsecase user.Usecase
}

func NewSubscriptionHandler(sUC subscriptions.Usecase,
	uUC user.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *SubscriptionHandler {
	return &SubscriptionHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		SUsecase: sUC,
		UUsecase: uUC,
	}
}

func (sh *SubscriptionHandler) Configure(e *echo.Echo) {
	e.GET("/api/v1/subscriptions", sh.GetSubscriptions(), sh.MManager.CheckAuth)
	e.POST("/api/v1/subscriptions", sh.Subscribe(), sh.MManager.CheckAuth)
	e.DELETE("/api/v1/subscriptions", sh.Unsubscribe(), sh.MManager.CheckAuth)
}

func (sh *SubscriptionHandler) GetSubscriptions() echo.HandlerFunc {
	type Request struct {
		Count  uint64 `query:"count" validate:"required"`
		Offset uint64 `query:"offset"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			sh.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := sh.ReqReader.Read(c, request, nil); err != nil {
			sh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		followers, totalFollowers, err := sh.UUsecase.GetFollowers(sess.UserID, request.Count, request.Offset)

		if err != nil {
			sh.Logger.Log(c, "error", "Error while getting followers.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		following, totalFollowing, err := sh.UUsecase.GetFollowing(sess.UserID, request.Count, request.Offset)

		if err != nil {
			sh.Logger.Log(c, "error", "Error while getting following.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"followers": followers,
				"total_followers": totalFollowers,
				"following": following,
				"total_following": totalFollowing,
			},
		})
	}
}

func (sh *SubscriptionHandler) Subscribe() echo.HandlerFunc {
	type Request struct {
		SubscriptionID uint64 `json:"subscription_id" validate:"required"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			sh.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := sh.ReqReader.Read(c, request, nil); err != nil {
			sh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := sh.SUsecase.Subscribe(sess.UserID, request.SubscriptionID); err != nil {
			sh.Logger.Log(c, "info", "Error while subscribing.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (sh *SubscriptionHandler) Unsubscribe() echo.HandlerFunc {
	type Request struct {
		SubscriptionID uint64 `json:"subscription_id" validate:"required"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			sh.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := sh.ReqReader.Read(c, request, nil); err != nil {
			sh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := sh.SUsecase.Unsubscribe(sess.ID, request.SubscriptionID); err != nil {
			sh.Logger.Log(c, "info", "Error while unsubscribing.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}
