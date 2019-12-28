package delivery

import (
	"2019_2_Covenant/pkg/likes"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/middlewares"
	"2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/reader"
	"2019_2_Covenant/pkg/user"
	. "2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	. "2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"net/http"
)

type LikesHandler struct {
	BaseHandler
	LUsecase likes.Usecase
	UUsecase user.Usecase
}

func NewLikesHandler(lUC likes.Usecase,
	uUC user.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *LikesHandler {
	return &LikesHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		LUsecase: lUC,
		UUsecase: uUC,
	}
}

func (lh *LikesHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/likes", lh.Like(), lh.MManager.CheckAuthStrictly)
	e.DELETE("/api/v1/likes", lh.Unlike(), lh.MManager.CheckAuthStrictly)
}

func (lh *LikesHandler) Like() echo.HandlerFunc {
	type Request struct {
		TrackID uint64 `json:"track_id" validate:"required"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			lh.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := lh.ReqReader.Read(c, request, nil); err != nil {
			lh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := lh.LUsecase.Like(sess.UserID, request.TrackID); err != nil {
			lh.Logger.Log(c, "info", "Error while putting like.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (lh *LikesHandler) Unlike() echo.HandlerFunc {
	type Request struct {
		TrackID uint64 `json:"track_id" validate:"required"`
	}

	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			lh.Logger.Log(c, "error", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := lh.ReqReader.Read(c, request, nil); err != nil {
			lh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := lh.LUsecase.Unlike(sess.UserID, request.TrackID); err != nil {
			lh.Logger.Log(c, "info", "Error while deleting like.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}
