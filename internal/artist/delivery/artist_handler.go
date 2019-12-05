package delivery

import (
	"2019_2_Covenant/internal/artist"
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
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

type ArtistHandler struct {
	BaseHandler
	AUsecase artist.Usecase
}

func NewArtistHandler(aUC artist.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *ArtistHandler {
	return &ArtistHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		AUsecase: aUC,
	}
}

func (ah *ArtistHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/artists", ah.CreateArtist(), ah.MManager.CheckAuth, ah.MManager.CheckAdmin)
	e.DELETE("/api/v1/artists/:id", ah.DeleteArtist(), ah.MManager.CheckAuth, ah.MManager.CheckAdmin)
	e.PUT("/api/v1/artists/:id", ah.UpdateArtist(), ah.MManager.CheckAuth, ah.MManager.CheckAdmin)
	//TODO: e.PUT("/api/v1/artists/:id/photo", ah.SetPhoto(), ah.MManager.CheckAuth, ah.MManager.CheckAdmin)
	e.GET("/api/v1/artists", ah.GetArtists())
	e.POST("/api/v1/artists/:id/albums", ah.CreateAlbum(), ah.MManager.CheckAuth, ah.MManager.CheckAdmin)
	//TODO: e.GET("/api/v1/artists/:id/albums", ah.GetAlbums())
}

func (ah *ArtistHandler) CreateAlbum() echo.HandlerFunc {
	type Request struct {
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

		a := models.NewAlbum(request.Name, request.Year, uint64(aID))

		if err := ah.AUsecase.CreateAlbum(a); err != nil {
			ah.Logger.Log(c, "info", "Error while storing album.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"album": a,
			},
		})
	}
}

func (ah *ArtistHandler) CreateArtist() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required"`
	}

	return func(c echo.Context) error {
		request := &Request{}

		if err := ah.ReqReader.Read(c, request, nil); err != nil {
			ah.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		a := &models.Artist{
			Name: request.Name,
		}

		if err := ah.AUsecase.Store(a); err != nil {
			ah.Logger.Log(c, "info", "Error while storing artist.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"artist": a,
			},
		})
	}
}

func (ah *ArtistHandler) DeleteArtist() echo.HandlerFunc {
	return func(c echo.Context) error {
		aID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ah.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		if err := ah.AUsecase.DeleteByID(uint64(aID)); err != nil {
			ah.Logger.Log(c, "info", "Error while deleting artist.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ah *ArtistHandler) UpdateArtist() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required"`
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

		if err := ah.ReqReader.Read(c, request, nil); err != nil {
			ah.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := ah.AUsecase.UpdateByID(uint64(aID), request.Name); err != nil {
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

func (ah *ArtistHandler) GetArtists() echo.HandlerFunc {
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

		artists, total, err := ah.AUsecase.Fetch(request.Count, request.Offset)

		if err != nil {
			ah.Logger.Log(c, "error", "Error while fetching artists", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"artists": artists,
				"total": total,
			},
		})
	}
}
