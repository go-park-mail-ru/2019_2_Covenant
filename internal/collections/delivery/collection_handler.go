package delivery

import (
	"2019_2_Covenant/internal/collections"
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
	"strconv"
)

type CollectionHandler struct {
	BaseHandler
	CUsecase collections.Usecase
}

func NewCollectionHandler(cUC collections.Usecase,
	mManager *middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *CollectionHandler {
	return &CollectionHandler{
		BaseHandler: BaseHandler{
			MManager:  mManager,
			Logger:    logger,
			ReqReader: reader.NewReqReader(),
		},
		CUsecase: cUC,
	}
}

func (ch *CollectionHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/collections", ch.CreateCollection(), ch.MManager.CheckAuthStrictly, ch.MManager.CheckAdmin)
	e.DELETE("/api/v1/collections/:id", ch.DeleteCollection(), ch.MManager.CheckAuthStrictly, ch.MManager.CheckAdmin)
	e.PUT("/api/v1/collections/:id", ch.UpdateCollection(), ch.MManager.CheckAuthStrictly, ch.MManager.CheckAdmin)
	e.GET("/api/v1/collections", ch.GetCollections())
	e.GET("/api/v1/collections/:id", ch.GetSingleCollection())
	e.POST("/api/v1/collections/:id/tracks", ch.AddToCollection(), ch.MManager.CheckAuthStrictly, ch.MManager.CheckAdmin)
	e.GET("/api/v1/collections/:id/tracks", ch.GetTracksFromCollection(), ch.MManager.CheckAuth)
	e.PUT("/api/v1/collections/:id/photo", ch.UploadCollectionPhoto(), ch.MManager.CheckAuthStrictly, ch.MManager.CheckAdmin)
}

func (ch *CollectionHandler) CreateCollection() echo.HandlerFunc {
	type Request struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	return func(c echo.Context) error {
		request := &Request{}

		if err := ch.ReqReader.Read(c, request, nil); err != nil {
			ch.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		collection := models.NewCollection(request.Name, request.Description)

		if err := ch.CUsecase.Store(collection); err != nil {
			ch.Logger.Log(c, "info", "Error while storing collection.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"collection": collection,
			},
		})
	}
}

func (ch *CollectionHandler) DeleteCollection() echo.HandlerFunc {
	return func(c echo.Context) error {
		cID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ch.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		if err := ch.CUsecase.DeleteByID(uint64(cID)); err != nil {
			ch.Logger.Log(c, "info", "Error while deleting collection.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ch *CollectionHandler) UpdateCollection() echo.HandlerFunc {
	type Request struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}

	return func(c echo.Context) error {
		cID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ch.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := ch.ReqReader.Read(c, request, nil); err != nil {
			ch.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := ch.CUsecase.UpdateByID(uint64(cID), request.Name, request.Description); err != nil {
			ch.Logger.Log(c, "info", "Error while updating collection.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ch *CollectionHandler) GetCollections() echo.HandlerFunc {
	type Request struct {
		Count  uint64 `query:"count" validate:"required"`
		Offset uint64 `query:"offset"`
	}

	return func(c echo.Context) error {
		request := &Request{}

		if err := ch.ReqReader.Read(c, request, nil); err != nil {
			ch.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		colls, total, err := ch.CUsecase.Fetch(request.Count, request.Offset)

		if err != nil {
			ch.Logger.Log(c, "error", "Error while fetching collections", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"collections": colls,
				"total":  total,
			},
		})
	}
}

func (ch *CollectionHandler) GetSingleCollection() echo.HandlerFunc {
	return func(c echo.Context) error {
		cID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ch.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		collection, amountOfTracks, err := ch.CUsecase.GetByID(uint64(cID))

		if err != nil {
			ch.Logger.Log(c, "info", "Error while getting collection", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"collection":       collection,
				"amount_of_tracks": amountOfTracks,
			},
		})
	}
}

func (ch *CollectionHandler) AddToCollection() echo.HandlerFunc {
	type Request struct {
		TrackID uint64 `json:"track_id" validate:"required"`
	}

	return func(c echo.Context) error {
		cID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ch.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		request := &Request{}

		if err := ch.ReqReader.Read(c, request, nil); err != nil {
			ch.Logger.Log(c, "info", "Invalid request.", err.Error())
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		if err := ch.CUsecase.AddTrack(uint64(cID), request.TrackID); err != nil {
			ch.Logger.Log(c, "info", "Error while adding track to collection.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ch *CollectionHandler) GetTracksFromCollection() echo.HandlerFunc {
	return func(c echo.Context) error {
		cID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ch.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		// TODO: Strictly + same handler in PlaylistHandler

		var authID uint64
		if sess, ok := c.Get("session").(*models.Session); ok {
			authID = sess.UserID
		}

		tracks, err := ch.CUsecase.GetTracks(uint64(cID), authID)

		if err != nil {
			ch.Logger.Log(c, "error", "Error while fetching tracks.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		for _, item := range tracks { item.Duration = time_parser.GetDuration(item.Duration) }

		return c.JSON(http.StatusOK, Response{
			Body: &Body{
				"tracks": tracks,
			},
		})
	}
}

func (ch *CollectionHandler) UploadCollectionPhoto() echo.HandlerFunc {
	rootPath, _ := os.Getwd()

	return func(c echo.Context) error {
		cID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			ch.Logger.Log(c, "error", "Atoi error.", err.Error())
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		file, err := c.FormFile("file")
		if err != nil {
			ch.Logger.Log(c, "info", "Can't extract file from request.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: ErrRetrievingError.Error(),
			})
		}

		src, err := file.Open()
		if err != nil {
			ch.Logger.Log(c, "error", "Can't open file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		defer src.Close()

		filePath := fmt.Sprintf("%s%s-%s", COLLECTIONS_PHOTOS_PATH, uuid.New().String(), file.Filename)

		dest, err := os.Create(filepath.Join(rootPath, filePath))
		if err != nil {
			ch.Logger.Log(c, "error", "Can't create file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		defer dest.Close()

		if _, err = io.Copy(dest, src); err != nil {
			ch.Logger.Log(c, "error", "Can't copy file.", err)
			return c.JSON(http.StatusInternalServerError, Response{
				Error: ErrInternalServerError.Error(),
			})
		}

		if err := ch.CUsecase.UpdatePhoto(uint64(cID), filePath); err != nil {
			ch.Logger.Log(c, "info", "Error while storing photo in db.", err)
			return c.JSON(http.StatusBadRequest, Response{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}
