package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/search"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	. "2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	"2019_2_Covenant/tools/time_parser"
	. "2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SearchHandler struct {
	BaseHandler
	SUsecase search.Usecase
	UUsecase user.Usecase
}

func NewSearchHandler(sUC search.Usecase, uUC user.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *SearchHandler {
	return &SearchHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		SUsecase: sUC,
		UUsecase: uUC,
	}
}

func (sh *SearchHandler) Configure(e *echo.Echo) {
	e.GET("/api/v1/search", sh.Search(), sh.MManager.CheckAuth)
}

func (sh *SearchHandler) Search() echo.HandlerFunc {
	type Request struct {
		Search  string `query:"s"`
	}

	isUserSearching := func(req *Request) bool {
		if req.Search[0] == '@' {
			req.Search = req.Search[1:]
			return true
		}
		return false
	}

	return func(c echo.Context) error {
		request := &Request{}

		if err := sh.ReqReader.Read(c, request, nil); err != nil {
			sh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		body := &Body{}

		if isUserSearching(request) {
			usr, err := sh.UUsecase.FindLike(request.Search, 10)
			if err != nil {
				sh.Logger.Log(c, "info", "Error while searching user.", err)
				return c.JSON(http.StatusNotFound, Response{
					Error: ErrNotFound.Error(),
				})
			}

			body = &Body{
				"user": usr,
			}
		} else {
			var authID uint64
			if sess, ok := c.Get("session").(*models.Session); ok {
				authID = sess.UserID
			}

			tracks, albums, artists, err := sh.SUsecase.Search(request.Search, 10, authID)

			if err != nil {
				sh.Logger.Log(c, "info", "Error while searching.", err)
				return c.JSON(http.StatusNotFound, Response{
					Error: ErrNotFound.Error(),
				})
			}

			for _, item := range tracks {
				item.Duration = time_parser.GetDuration(item.Duration)
			}

			body = &Body{
				"tracks": tracks,
				"albums": albums,
				"artists": artists,
			}
		}

		return c.JSON(http.StatusOK, Response{
			Body: body,
		})
	}
}
