package handlers

import (
	. "../storage"
	"encoding/json"
	"net/http"
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
		http.Error(w, `{"error":"database"}`, 500)
		return
	}

	api.SignIn(w, r)
}

func (api *UsersHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	id, err := api.store.CheckUser(email, password)
	if err != nil {
		http.Error(w, `{"error":"email and password are mismatched"}`, 500)
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

	err = json.NewEncoder(w).Encode(&Result{Body: body})
	if err != nil {
		http.Error(w, `{"error":"wrong data"}`, 500)
		return
	}
}
