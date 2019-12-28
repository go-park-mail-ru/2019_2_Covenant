package delivery

import (
	"2019_2_Covenant/pkg/album"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/middlewares"
	"2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/reader"
	. "2019_2_Covenant/tools/base_handler"
	. "2019_2_Covenant/tools/response"
	"2019_2_Covenant/tools/time_parser"
	. "2019_2_Covenant/tools/vars"
	"encoding/json"
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
	e.DELETE("/api/v1/albums/:id", ah.DeleteAlbum(), ah.MManager.CheckAuthStrictly, ah.MManager.CheckAdmin)
	e.PUT("/api/v1/albums/:id", ah.UpdateAlbum(), ah.MManager.CheckAuthStrictly, ah.MManager.CheckAdmin)
	e.GET("/api/v1/albums", ah.GetAlbums())
	e.GET("/api/v1/albums/:id", ah.GetSingleAlbum())
	e.POST("/api/v1/albums/:id/tracks", ah.AddToAlbum(), ah.MManager.CheckAuthStrictly, ah.MManager.CheckAdmin)
	e.GET("/api/v1/albums/:id/tracks", ah.GetTracksFromAlbum(), ah.MManager.CheckAuth)
	e.PUT("/api/v1/albums/:id/photo", ah.UploadAlbumPhoto(), ah.MManager.CheckAuthStrictly, ah.MManager.CheckAdmin)
}

func (ah *AlbumHandler) UploadAlbumPhoto() echo.HandlerFunc {
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

		if err := ah.AUsecase.UpdatePhoto(c.Request().Context(), uint64(aID), src); err != nil {
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ah *AlbumHandler) GetTracksFromAlbum() echo.HandlerFunc {
	return func(c echo.Context) error {
		aID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ah.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		var authID uint64
		if sess, ok := c.Get("session").(*models.Session); ok {
			authID = sess.UserID
		}

		tracks, err := ah.AUsecase.GetTracksFrom(uint64(aID), authID)

		if err != nil {
			ah.Logger.Log(c, "error", "Error while fetching tracks.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		for _, item := range tracks {
			item.Duration = time_parser.GetDuration(item.Duration)
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"tracks": tracks,
			},
		})
	}
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

func (ah *AlbumHandler) GetSingleAlbum() echo.HandlerFunc {
	return func(c echo.Context) error {
		aID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ah.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		a, amountOfTracks, err := ah.AUsecase.GetByID(uint64(aID))

		if err != nil {
			ah.Logger.Log(c, "info", "Error while getting album", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"album":            a,
				"amount_of_tracks": amountOfTracks,
			},
		})
	}
}

func (ah *AlbumHandler) AddToAlbum() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required"`
	}

	return func(c echo.Context) error {
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

		aID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ah.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		form, _ := c.MultipartForm()
		request := &Request{}

		if err := json.Unmarshal([]byte(form.Value["request"][0]), request); err != nil {
			ah.Logger.Log(c, "info", "Error while parsing JSON.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := ah.ReqReader.Read(c, request, nil); err != nil {
			ah.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := ah.AUsecase.AddTrack(c.Request().Context(), uint64(aID), request.Name, src); err != nil {
			ah.Logger.Log(c, "error", "Error while adding track to album.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrAlreadyExist.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}
