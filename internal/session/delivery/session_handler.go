package delivery

import (
	"2019_2_Covenant/internal/middlewares"
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/internal/vars"
	"2019_2_Covenant/pkg/logger"
	"2019_2_Covenant/pkg/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

type SessionHandler struct {
	SUsecase     session.Usecase
	UUsecase     user.Usecase
	MManager     middlewares.MiddlewareManager
	Logger       *logger.LogrusLogger
	ReqValidator *validator.ReqValidator
}

func NewSessionHandler(sUC session.Usecase,
	uUC user.Usecase,
	mManager middlewares.MiddlewareManager,
	logger *logger.LogrusLogger) *SessionHandler {
	return &SessionHandler{
		SUsecase:     sUC,
		UUsecase:     uUC,
		MManager:     mManager,
		Logger:       logger,
		ReqValidator: validator.NewReqValidator(),
	}
}

func (sh *SessionHandler) Configure(e *echo.Echo) {
	e.POST("/api/v1/signup", sh.SignUp())
	e.POST("/api/v1/login", sh.LogIn())
	e.GET("/api/v1/logout", sh.LogOut(), sh.MManager.CheckAuth)
	e.GET("/api/v1/csrf", sh.GetCSRF(), sh.MManager.CheckAuth)
}

// @Tags User
// @Summary SignUp Route
// @Description Signing user up
// @ID sign-up-user
// @Accept json
// @Produce json
// @Param Data body object true "JSON that contains user sign up data"
// @Success 200 object models.User
// @Failure 400 object vars.ResponseError
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
// @Router /api/v1/signup [post]
func (sh *SessionHandler) SignUp() echo.HandlerFunc {
	type UserReg struct {
		Nickname         string `json:"nickname" validate:"required"`
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,gte=6"`
		PassConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	}

	correctData := func(data *UserReg) bool {
		return strings.Contains(data.Password, " ") == false &&
			strings.Contains(data.Nickname, " ") == false
	}

	return func(c echo.Context) error {
		var userRegData UserReg
		err := c.Bind(&userRegData)

		if err != nil {
			sh.Logger.Log(c, "error", "Can't read request body.")
			return c.JSON(http.StatusUnprocessableEntity, vars.ResponseError{Error: err.Error()})
		}

		if err := sh.ReqValidator.Validate(userRegData); err != nil || !correctData(&userRegData) {
			sh.Logger.Log(c, "info", "Invalid request.", userRegData)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
		}

		usr, err := sh.UUsecase.GetByEmail(userRegData.Email)

		if usr != nil {
			sh.Logger.Log(c, "info", "Already exist.", "User ID:", usr.ID)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{
				Error: vars.ErrAlreadyExist.Error(),
			})
		}

		newUser := &models.User{
			Email:         userRegData.Email,
			PlainPassword: userRegData.Password,
			Nickname:      userRegData.Nickname,
		}

		usr, err = sh.UUsecase.Store(newUser)

		if err != nil {
			sh.Logger.Log(c, "error", "User store error.", err)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
		}

		cookie := &http.Cookie{
			Name:    "Covenant",
			Value:   uuid.New().String(),
			Expires: time.Now().Add(24 * time.Hour),
		}

		sess := &models.Session{
			UserID:  usr.ID,
			Expires: cookie.Expires,
			Data:    cookie.Value,
		}

		err = sh.SUsecase.Store(sess)

		if err != nil {
			sh.Logger.Log(c, "error", "Session store error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		token, err := models.NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, time.Now().Add(24*time.Hour))
		c.Response().Header().Set("X-CSRF-Token", token)

		if err != nil {
			sh.Logger.Log(c, "error", "CSRF Token generating error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, vars.Response{Body: newUser})
	}
}

// @Tags User
// @Summary LogIn Route
// @Description Logging user in
// @ID log-in-user
// @Accept json
// @Produce json
// @Param Data body object true "JSON that contains user login data"
// @Success 200 object models.User
// @Failure 400 object vars.ResponseError
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
// @Router /api/v1/login [post]
func (sh *SessionHandler) LogIn() echo.HandlerFunc {
	type UserLogin struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	return func(c echo.Context) error {
		var userLoginData UserLogin
		err := c.Bind(&userLoginData)

		if err != nil {
			sh.Logger.Log(c, "error", "Can't read request body.")
			return c.JSON(http.StatusUnprocessableEntity, vars.ResponseError{Error: err.Error()})
		}

		if err := sh.ReqValidator.Validate(userLoginData); err != nil {
			sh.Logger.Log(c, "info", "Invalid request.", userLoginData)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
		}

		usr, err := sh.UUsecase.GetByEmail(userLoginData.Email)

		if usr == nil {
			sh.Logger.Log(c, "info", "Error while getting user by EMAIL.", err)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{Error: err.Error()})
		}

		if !usr.Verify(userLoginData.Password) {
			sh.Logger.Log(c, "info", "Bad authentication.", "User:", usr.Nickname)
			return c.JSON(http.StatusBadRequest, vars.ResponseError{
				Error: vars.ErrBadParam.Error(),
			})
		}

		cookie := &http.Cookie{
			Name:    "Covenant",
			Value:   uuid.New().String(),
			Expires: time.Now().Add(24 * time.Hour),
		}

		sess := &models.Session{
			UserID:  usr.ID,
			Expires: cookie.Expires,
			Data:    cookie.Value,
		}

		err = sh.SUsecase.Store(sess)

		if err != nil {
			sh.Logger.Log(c, "error", "Session store error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		token, err := models.NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, time.Now().Add(24*time.Hour))
		c.Response().Header().Set("X-CSRF-Token", token)

		if err != nil {
			sh.Logger.Log(c, "error", "CSRF Token generating error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, vars.Response{Body: usr})
	}
}

// @Tags User
// @Summary Log Out Route
// @Description Logging user out
// @ID log-out-user
// @Accept json
// @Produce json
// @Success 200 object models.User
// @Failure 400 object vars.ResponseError
// @Failure 404 object vars.ResponseError
// @Failure 500 object vars.ResponseError
// @Router /api/v1/logout [get]
func (sh *SessionHandler) LogOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			sh.Logger.Log(c, "info", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		if err := sh.SUsecase.DeleteByID(sess.ID); err != nil {
			sh.Logger.Log(c, "error", "Error while deleting session.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		cookie := &http.Cookie{
			Name:    "Covenant",
			Value:   sess.Data,
			Expires: time.Now().AddDate(0, 0, -1),
		}

		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]string{
			"message": "logout",
		})
	}
}

func (sh *SessionHandler) GetCSRF() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, ok := c.Get("session").(*models.Session)

		if !ok {
			sh.Logger.Log(c, "info", "Can't extract session from echo.Context.")
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: vars.ErrInternalServerError.Error(),
			})
		}

		token, err := models.NewCSRFTokenManager("Covenant").Create(sess.UserID, sess.Data, time.Now().Add(24*time.Hour))
		c.Response().Header().Set("X-CSRF-Token", token)

		if err != nil {
			sh.Logger.Log(c, "error", "CSRF Token generating error.", err)
			return c.JSON(http.StatusInternalServerError, vars.ResponseError{
				Error: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "success",
		})
	}
}
