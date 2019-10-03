package handlers

import (
	. "../storage"
	"encoding/json"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"net/http"
	"time"
)

type Result struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
}

type UsersHandler struct {
	store *UserStore
	session *SessionsStore
}

func (api *UsersHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	avatar := r.FormValue("avatar")

	if avatar == "" {
		avatar = "img/user_profile.png"
	}

	newUser := &User{
		Username: username,
		Email: email,
		Password: password,
		Avatar: avatar,
	}

	_, err := api.store.AddUser(newUser)
	if err != nil {
		Error(w, err, 500)
		return
	}

	api.SignIn(w, r)
}

func (api *UsersHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	id, err := api.store.CheckUser(email, password)
	if err != nil {
		Error(w, err, 500)
		return
	}

	sessionID, session := api.session.Set(id)

	cookie := http.Cookie{
		Name:    "session-id",
		Value:   sessionID,
		Expires: session.Expires,
	}

	http.SetCookie(w, &cookie)

	body := map[string]interface{}{
		"id": id,
	}

	_ = json.NewEncoder(w).Encode(&Result{Body: body})
}

func (api *UsersHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err != nil {
		Error(w, err, 401)
	}

	_, err = api.session.Get(session.Value)
	if err != nil {
		Error(w, err, 401)
		return
	}

	api.session.Delete(session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func Error(w http.ResponseWriter, error error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	_ := json.NewEncoder(w).Encode(&Result{Err: fmt.Sprint(error)})
}