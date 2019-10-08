package handlers

import (
	. "../storage"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	Body interface{} `json:"body,omitempty"`
	Err  string      `json:"err,omitempty"`
	Auth bool        `json:"auth"`
}

type UsersHandler struct {
	store *UserStore
	session *SessionsStore
}

func (api *UsersHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	if api.store.IsExist(email) {
		err := fmt.Errorf("user exists")
		Error(w, err, 500)
		return
	}

	newUser := &User{
		Email: email,
		Password: password,
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

	_ = json.NewEncoder(w).Encode(&Result{Body: body, Auth: true})
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

	_ = json.NewEncoder(w).Encode(&Result{Auth: false})
}

func Error(w http.ResponseWriter, error error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	_ := json.NewEncoder(w).Encode(&Result{Err: fmt.Sprint(error)})
}