package usecase

import (
	. "2019_2_Covenant/internal/models"
	user2 "2019_2_Covenant/internal/user"
	vars "2019_2_Covenant/internal/vars"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Users struct {
	User []*User
}

var users = Users{
	User: []*User{
		{ID: 1, Username: "marshal", Email: "m1@ya.ru", Password: "12345"},
		{ID: 2, Username: "plaksenka", Email: "p2@ya.ru", Password: "12345"},
		{ID: 3, Username: "svya", Email: "s3@ya.ru", Password: "12345"},
	},
}

func TestUserUsecase_FetchAllOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := user2.NewMockRepository(ctrl)
	userRepo.EXPECT().FetchAll().Return(users.User, nil)
	usecase := userUsecase{userRepo: userRepo}

	expUsers, err := usecase.FetchAll()

	asserts := assert.New(t)

	asserts.Equal(users.User, expUsers, "Should be equal")
	asserts.Nil(err)
}


func TestUserUsecase_FetchAllFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := user2.NewMockRepository(ctrl)
	userRepo.EXPECT().FetchAll().Return(nil, vars.ErrNotFound)
	usecase := userUsecase{userRepo: userRepo}

	expUsers, err := usecase.FetchAll()

	asserts := assert.New(t)

	asserts.Nil(expUsers)
	asserts.NotNil(err)
}