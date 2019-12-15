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
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	e.POST("/api/v1/artists", ah.CreateArtist(), ah.MManager.CheckAuthStrictly, ah.MManager.CheckAdmin)
	e.DELETE("/api/v1/artists/:id", ah.DeleteArtist(), ah.MManager.CheckAuthStrictly, ah.MManager.CheckAdmin)
	e.PUT("/api/v1/artists/:id", ah.UpdateArtist(), ah.MManager.CheckAuthStrictly, ah.MManager.CheckAdmin)
	e.PUT("/api/v1/artists/:id/photo", ah.UploadArtistPhoto(), ah.MManager.CheckAuthStrictly, ah.MManager.CheckAdmin)
	e.GET("/api/v1/artists", ah.GetArtists())
	e.GET("/api/v1/artists/:id", ah.GetSingleArtist())
	e.POST("/api/v1/artists/:id/albums", ah.CreateAlbum(), ah.MManager.CheckAuthStrictly, ah.MManager.CheckAdmin)
	e.GET("/api/v1/artists/:id/albums", ah.GetArtistAlbums())
	e.GET("/api/v1/artists/:id/tracks", ah.GetArtistTracks())
}

func (ah *ArtistHandler) GetArtistTracks() echo.HandlerFunc {
	type Request struct {
		Count  uint64 `query:"count" validate:"required"`
		Offset uint64 `query:"offset"`
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

		tracks, total, err := ah.AUsecase.GetTracks(uint64(aID), request.Count, request.Offset)

		if err != nil {
			ah.Logger.Log(c, "error", "Error while getting artist tracks.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"tracks": tracks,
				"total": total,
			},
		})
	}
}

func (ah *ArtistHandler) GetArtistAlbums() echo.HandlerFunc {
	type Request struct {
		Count  uint64 `query:"count" validate:"required"`
		Offset uint64 `query:"offset"`
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

		albums, total, err := ah.AUsecase.GetArtistAlbums(uint64(aID), request.Count, request.Offset)

		if err != nil {
			ah.Logger.Log(c, "error", "Error while fetching artist's albums", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"albums": albums,
				"total": total,
			},
		})
	}
}

func (ah *ArtistHandler) UploadArtistPhoto() echo.HandlerFunc {
	rootPath, _ := os.Getwd()

	return func(c echo.Context) error {
		aID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ah.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		file, err := c.FormFile("file")
		if err != nil {
			ah.Logger.Log(c, "info", "Can't extract file from request.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: ErrRetrievingError.Error(),
			})
		}

		src, err := file.Open()
		if err != nil {
			ah.Logger.Log(c, "error", "Can't open file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		defer src.Close()

		filePath := fmt.Sprintf("%s%s-%s", ARTISTS_PHOTOS_PATH, uuid.New().String(), file.Filename)

		dest, err := os.Create(filepath.Join(rootPath, filePath))
		if err != nil {
			ah.Logger.Log(c, "error", "Can't create file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		defer dest.Close()

		if _, err = io.Copy(dest, src); err != nil {
			ah.Logger.Log(c, "error", "Can't copy file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		if err := ah.AUsecase.UpdatePhoto(uint64(aID), filePath); err != nil {
			ah.Logger.Log(c, "info", "Error while storing photo in db.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
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

func (ah *ArtistHandler) GetSingleArtist() echo.HandlerFunc {
	return func(c echo.Context) error {
		pID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ah.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		a, amountOfAlbums, err := ah.AUsecase.GetByID(uint64(pID))

		if err != nil {
			ah.Logger.Log(c, "info", "Error while getting playlist.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"artist": a,
				"amount_of_albums": amountOfAlbums,
			},
		})
	}
}
