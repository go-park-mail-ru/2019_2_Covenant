package usecase

import (
	. "2019_2_Covenant/internal/models"
	mock "2019_2_Covenant/internal/user/mocks"
	"2019_2_Covenant/internal/vars"
	"github.com/golang/mock/gomock"
	"testing"
)

//go:generate mockgen -source=../repository.go -destination=../mocks/mock_repository.go -package=mock

type Users struct {
	User []*User
}

var users = Users{
	User: []*User{
		{ID: 1, Username: "marshal", Email: "m1@ya.ru", Password: "123456"},
		{ID: 2, Username: "plaksenka", Email: "p2@ya.ru", Password: "123456"},
		{ID: 3, Username: "svya", Email: "s3@ya.ru", Password: "123456"},
	},
}

func TestUserUsecase_FetchAllOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)
	userRepo.EXPECT().FetchAll().Return(users.User, nil)
	usecase := userUsecase{userRepo: userRepo}

	expUsers, err := usecase.FetchAll()

	if gomock.Not(users.User).Matches(expUsers) || err != nil {
		t.Fail()
	}
}


func TestUserUsecase_FetchAllFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)
	userRepo.EXPECT().FetchAll().Return(nil, vars.ErrNotFound)
	usecase := userUsecase{userRepo: userRepo}

	expUsers, err := usecase.FetchAll()

	if expUsers != nil || err == nil {
		t.Fail()
	}
}
