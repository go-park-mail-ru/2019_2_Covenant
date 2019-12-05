package delivery

import (
	"2019_2_Covenant/internal/album"
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/reader"
	. "2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	"2019_2_Covenant/tools/time_parser"
	. "2019_2_Covenant/tools/vars"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type AlbumHandler struct {
	BaseHandler
	AUsecase album.Usecase
}

func NewAlbumHandler(aUC album.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *AlbumHandler {
	return &AlbumHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		AUsecase: aUC,
	}
}

func (ah *AlbumHandler) Configure(e *echo.Echo) {
	e.DELETE("/api/v1/albums/:id", ah.DeleteAlbum(), ah.MManager.CheckAuth, ah.MManager.CheckAdmin)
	e.PUT("/api/v1/albums/:id", ah.UpdateAlbum(), ah.MManager.CheckAuth, ah.MManager.CheckAdmin)
	e.GET("/api/v1/albums", ah.GetAlbums())
}

func (ah *AlbumHandler) DeleteAlbum() echo.HandlerFunc {
	return func(c echo.Context) error {
		aID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ah.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		if err := ah.AUsecase.DeleteByID(uint64(aID)); err != nil {
			ah.Logger.Log(c, "info", "Error while deleting album.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ah *AlbumHandler) UpdateAlbum() echo.HandlerFunc {
	type Request struct {
		ArtistID uint64 `json:"artist_id" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Year     string `json:"year" validate:"required"`
	}

	correctData := func(req interface{}) bool {
		reg, err := regexp.Compile("^[0-9-]*$")

		if err != nil || !reg.MatchString(req.(*Request).Year) {
			return false
		}

		timeNow := time.Now()
		date := time_parser.StringToTime(req.(*Request).Year)

		if date.Sub(timeNow) > 0 {
			return false
		}

		return true
	}

	return func(c echo.Context) error {
		aID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ah.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := ah.ReqReader.Read(c, request, correctData); err != nil {
			ah.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := ah.AUsecase.UpdateByID(uint64(aID), request.ArtistID, request.Name, request.Year); err != nil {
			ah.Logger.Log(c, "info", "Error while updating artist.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ah *AlbumHandler) GetAlbums() echo.HandlerFunc {
	type Request struct {
		Count  uint64 `query:"count" validate:"required"`
		Offset uint64 `query:"offset"`
	}

	return func(c echo.Context) error {
		request := &Request{}

		if err := ah.ReqReader.Read(c, request, nil); err != nil {
			ah.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		albums, total, err := ah.AUsecase.Fetch(request.Count, request.Offset)

		if err != nil {
			ah.Logger.Log(c, "error", "Error while fetching artists", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"albums": albums,
				"total":  total,
			},
		})
	}
}
