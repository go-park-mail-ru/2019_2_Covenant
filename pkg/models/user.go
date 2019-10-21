package models

type User struct {
	ID			uint64   `json:"-"`
	Username	string 	`json:"username"`
	Email		string 	`json:"email"`
	Name		string 	`json:"name"`
	Surname		string 	`json:"surname"`
	Password	string 	`json:"-"`
	Avatar		string 	`json:"-"`
	Role		int8   	`json:"role"` // 0 - user; 1 - admin;
	Access		int8   	`json:"access"` // 0 - public; 1 - private;
}

type UserReg struct {
	Username	string   `json:"username" validate:"required"`
	Email		string   `json:"email" validate:"required,email"`
	Password	string   `json:"password" validate:"required,gte=6"`
}

type UserLogin struct {
	Email		string   `json:"email" validate:"required,email"`
	Password	string   `json:"password" validate:"required"`
}

type UserEdit struct {
	Name		string   `json:"name" validate:"required"`
	Surname		string   `json:"surname" validate:"required"`
}
