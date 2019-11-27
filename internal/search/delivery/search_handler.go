package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/search"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	. "2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	. "2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SearchHandler struct {
	BaseHandler
	SUsecase search.Usecase
}

func NewSearchHandler(sUC search.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *SearchHandler {
	return &SearchHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		SUsecase: sUC,
	}
}

func (sh *SearchHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/search", sh.Search())
}

func (sh *SearchHandler) Search() echo.HandlerFunc {
	type Request struct {
		Text  string `json:"text"`
		Count uint64 `json:"count"`
	}

	return func(c echo.Context) error {
		request := &Request{}

		if err := sh.ReqReader.Read(c, request, nil); err != nil {
			sh.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		tracks, albums, artists, err := sh.SUsecase.Search(request.Text, request.Count)

		if err != nil {
			sh.Logger.Log(c, "info", "Error while searching.", err)
			return c.JSON(http.StatusNotFound, Response{
				Error: ErrNotFound.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"tracks":  tracks,
				"albums":  albums,
				"artists": artists,
			},
		})
	}
}
