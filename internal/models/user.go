package models

type User struct {
	ID       uint64 `json:"-"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"-"`
	Avatar   string `json:"avatar"`
	Role     int8   `json:"role"`   // 0 - user; 1 - admin;
	Access   int8   `json:"access"` // 0 - public; 1 - private;
}
