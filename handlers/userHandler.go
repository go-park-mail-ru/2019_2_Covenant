package handlers

import (
	. "../storage"
	"encoding/json"
	"fmt"
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
	Store   *UserStore
	Session *SessionsStore
}

func (api *UsersHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	if api.Store.IsExist(email) {
		err := fmt.Errorf("user exists")
		Error(w, false, err, 500)
		return
	}

	newUser := &User{
		Email: email,
		Password: password,
	}

	_, err := api.Store.AddUser(newUser)
	if err != nil {
		Error(w, false, err, 500)
		return
	}

	api.SignIn(w, r)
}

func (api *UsersHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	auth := false

	email := r.FormValue("email")
	password := r.FormValue("password")

	id, err := api.Store.CheckUser(email, password)
	if err != nil {
		Error(w, auth, err, 500)
		return
	}

	user, err := api.Store.GetUserByID(id)
	if err != nil {
		Error(w, auth, err, 500)
		return
	}

	sessionID, session := api.Session.Set(id)

	cookie := http.Cookie{
		Name:    "session-id",
		Value:   sessionID,
		Expires: session.Expires,
	}

	http.SetCookie(w, &cookie)

	auth = true

	body := map[string]interface{}{
		"id": id,
		"username": user.Username,
		"avatar": user.Avatar,
	}

	_ = json.NewEncoder(w).Encode(&Result{Body: body, Auth: auth})
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