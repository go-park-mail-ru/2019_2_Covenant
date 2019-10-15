package handlers

import (
	"2019_2_Covenant/storage"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type Result struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
	Auth bool        `json:"auth"`
}

type UsersHandler struct {
	Store   *storage.UserStore
	Session *storage.SessionsStore
}

type Body map[string]interface{}

func (b Body) toString(key string) string {
	return fmt.Sprint(b[key])
}

func (b Body) contain(key string) bool {
	_, exist := b[key]
	return exist
}

func (api *UsersHandler) SignUp(c echo.Context) error {
	body := make(Body)
	err := json.NewDecoder(c.Request().Body).Decode(&body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
	}

	if !(body.contain("email") && body.contain("password")) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "bad params",
		})
	}

	if api.Store.IsExist(body.toString("email")) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "user exists",
		})
	}

	newUser := &User{
		Email: body.toString("email"),
		Password: body.toString("password"),
	}

	_, _ = api.Store.AddUser(newUser)

	cookie := new(http.Cookie)
	cookie.Name = "Covenant"
	cookie.Value = uuid.New().String()
	cookie.Expires = time.Now().Add(24 * time.Hour)

	session := &Session{
		UserID:  newUser.ID,
		Expires: cookie.Expires,
		Data:    cookie.Value,
	}

	_, _ = api.Session.Set(session)

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, session)
}

func (api *UsersHandler) SignIn(c echo.Context) error {
	body := make(Body)
	err := json.NewDecoder(c.Request().Body).Decode(&body)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
	}

	if !(body.contain("email") && body.contain("password")) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "bad params",
		})
	}

	if !api.Store.IsExist(body.toString("email")) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "user does not exist",
		})
	}

	user, _ := api.Store.GetUserByEmail(body.toString("email"))

	if body.toString("password") != user.Password {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "wrong password",
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "Covenant"
	cookie.Value = uuid.New().String()
	cookie.Expires = time.Now().Add(24 * time.Hour)

	session := &Session{
		UserID:  user.ID,
		Expires: cookie.Expires,
		Data:    cookie.Value,
	}

	_, _ = api.Session.Set(session)

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, cookie.Value)
}

func (api *UsersHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	auth := false
	cookie, err := r.Cookie("session_id")
	if err != nil {
		Error(w, auth, err, 401)
		return
	}

	_, err = api.Session.Get(cookie.Value)
	if err != nil {
		Error(w, auth, err, 401)
		return
	}

	api.Session.Delete(cookie.Value)

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)

	_ = json.NewEncoder(w).Encode(&Result{Auth: auth})
}

func (api *UsersHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	auth := false

	cookie, err := r.Cookie("session_id")
	if err != nil {
		Error(w, auth, err, 401)
		return
	}

	_, err = api.Session.Get(cookie.Value)
	if err != nil {
		Error(w, auth, err, 401)
		return
	}

	auth = true

	id := r.FormValue("id")

	profileID, err := strToUint(id)
	if err != nil {
		Error(w, auth, err, 500)
		return
	}

	user, err := api.Store.GetUserByID(profileID)
	if err != nil {
		Error(w, auth, err, 500)
		return
	}

	body := map[string]interface{}{
		"id": user.ID,
		"username": user.Username,
		"email": user.Email,
		"avatar": user.Avatar,
	}

	_ = json.NewEncoder(w).Encode(&Result{Body: body, Auth: auth})
}

func (api *UsersHandler) PostProfile(w http.ResponseWriter, r *http.Request) {

	auth := false

	// Проверяем куки
	cookie, err := r.Cookie("session_id")
	if err != nil {
		Error(w, auth, err, 401)
		return
	}

	session, err := api.Session.Get(cookie.Value)
	if err != nil {
		Error(w, auth, err, 401)
		return
	}

	auth = true

	id := r.FormValue("id")
	newUsername := r.FormValue("username")

	profileID, err := strToUint(id)
	if err != nil {
		Error(w, auth, err, 500)
		return
	}

	// Вовзращаем ошибку, если пользователь пытается изменить не свои данные
	if profileID != session.UserID {
		err := fmt.Errorf("access error")
		Error(w, auth, err, 500)
		return
	}

	newUser, err := api.Store.ChangeUsername(profileID, newUsername)
	if err != nil {
		Error(w, auth, err, 500)
		return
	}

	body := map[string]interface{}{
		"id": newUser.ID,
		"username": newUser.Username,
		"email": newUser.Email,
		"avatar": newUser.Avatar,
	}

	_ = json.NewEncoder(w).Encode(&Result{Body: body, Auth: auth})
}


func Error(w http.ResponseWriter, auth bool, error error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(&Result{Err: fmt.Sprint(error), Auth: auth})
}

func strToUint(str string) (res uint, err error) {

	data, errParse := strconv.ParseUint(str, 10, 0)
	if errParse != nil {
		err = fmt.Errorf("data error")
		return
	}

	res = uint(data)

	return
}