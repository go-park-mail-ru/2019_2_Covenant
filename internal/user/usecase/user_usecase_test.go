package usecase

import (
	. "2019_2_Covenant/internal/models"
	mock "2019_2_Covenant/internal/user/mocks"
	"2019_2_Covenant/internal/vars"
	"fmt"
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


func TestUserUsecase_FetchAllErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)
	userRepo.EXPECT().FetchAll().Return(nil, fmt.Errorf("error"))
	usecase := userUsecase{userRepo: userRepo}

	expUsers, err := usecase.FetchAll()

	if expUsers != nil || err == nil {
		t.Fail()
	}
}

func TestUserUsecase_GetByEmailOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)
	userRepo.EXPECT().GetByEmail("m1@ya.ru").Return(users.User[0], nil)
	usecase := userUsecase{userRepo: userRepo}

	expUser, err := usecase.GetByEmail("m1@ya.ru")

	if expUser != users.User[0] || err != nil {
		t.Fail()
	}
}

func TestUserUsecase_GetByEmailErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)
	userRepo.EXPECT().GetByEmail("notfound@ya.ru").Return(nil, vars.ErrNotFound)
	usecase := userUsecase{userRepo: userRepo}

	expUser, err := usecase.GetByEmail("notfound@ya.ru")

	if expUser != nil || err != vars.ErrNotFound {
		t.Fail()
	}
}

func TestUserUsecase_GetByIDOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)
	userRepo.EXPECT().GetByID(uint64(1)).Return(users.User[0], nil)
	usecase := userUsecase{userRepo: userRepo}

	expUser, err := usecase.GetByID(1)

	if expUser != users.User[0] || err != nil {
		t.Fail()
	}
}

func TestUserUsecase_GetByIDErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)
	userRepo.EXPECT().GetByID(uint64(5)).Return(nil, vars.ErrNotFound)
	usecase := userUsecase{userRepo: userRepo}

	expUser, err := usecase.GetByID(5)

	if expUser != nil || err != vars.ErrNotFound {
		t.Fail()
	}
}
func TestUserUsecase_StoreOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newUser := User{Username: "newUser", Email: "n4@ya.ru", Password: "123456"}

	userRepo := mock.NewMockRepository(ctrl)
	userRepo.EXPECT().GetByEmail("n4@ya.ru").Return(nil, vars.ErrNotFound)
	userRepo.EXPECT().Store(&newUser).Return(nil)

	usecase := userUsecase{userRepo: userRepo}

	err := usecase.Store(&newUser)

	if err != nil {
		t.Fail()
	}
}

func TestUserUsecase_StoreErr1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newUser := User{Username: "newUser", Email: "m1@ya.ru", Password: "123456"}

	userRepo := mock.NewMockRepository(ctrl)
	userRepo.EXPECT().GetByEmail("m1@ya.ru").Return(users.User[0], nil)

	usecase := userUsecase{userRepo: userRepo}

	err := usecase.Store(&newUser)

	if err != vars.ErrAlreadyExist {
		t.Fail()
	}
}

func TestUserUsecase_StoreErr2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newUser := User{Username: "newUser", Email: "n4@ya.ru", Password: "123456"}

	userRepo := mock.NewMockRepository(ctrl)
	userRepo.EXPECT().GetByEmail("n4@ya.ru").Return(nil, vars.ErrNotFound)
	userRepo.EXPECT().Store(&newUser).Return(fmt.Errorf("some error"))

	usecase := userUsecase{userRepo: userRepo}

	err := usecase.Store(&newUser)

	if err != vars.ErrInternalServerError {
		t.Fail()
	}
}