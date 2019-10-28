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

func configUserUsecase(userRepo *mock.MockRepository) userUsecase {
	return userUsecase{userRepo: userRepo}
}

func TestUserUsecase_FetchAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func (usecase userUsecase) ([]*User, error) {
		return usecase.FetchAll()
	}

	t.Run("Test OK", func (t1 *testing.T) {
		userRepo.EXPECT().FetchAll().Return(users.User, nil)
		usecase := configUserUsecase(userRepo)

		expUsers, err := exe(usecase)

		if gomock.Not(users.User).Matches(expUsers) || err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func (t2 *testing.T) {
		userRepo.EXPECT().FetchAll().Return(nil, fmt.Errorf("error"))
		usecase := configUserUsecase(userRepo)

		expUsers, err := exe(usecase)

		if expUsers != nil || err == nil {
			t2.Fail()
		}
	})

}

func TestUserUsecase_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func (usecase userUsecase, email string) (*User, error) {
		return usecase.GetByEmail(email)
	}

	t.Run("Test OK", func (t1 *testing.T) {
		setEmail := "m1@ya.ru"

		userRepo.EXPECT().GetByEmail(setEmail).Return(users.User[0], nil)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, setEmail)

		if expUser != users.User[0] || err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func (t2 *testing.T) {
		setEmail := "notfound@ya.ru"

		userRepo.EXPECT().GetByEmail(setEmail).Return(nil, vars.ErrNotFound)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, setEmail)

		if expUser != nil || err != vars.ErrNotFound {
			t2.Fail()
		}
	})
}

func TestUserUsecase_GetByIDOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func (usecase userUsecase, ID uint64) (*User, error) {
		return usecase.GetByID(ID)
	}

	t.Run("Test OK", func (t1 *testing.T) {
		setID := uint64(1)

		userRepo.EXPECT().GetByID(setID).Return(users.User[0], nil)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, setID)

		if expUser != users.User[0] || err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func (t2 *testing.T) {
		setID := uint64(5)
		userRepo.EXPECT().GetByID(setID).Return(nil, vars.ErrNotFound)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, setID)

		if expUser != nil || err != vars.ErrNotFound {
			t2.Fail()
		}
	})

}

func TestUserUsecase_StoreOk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()



	userRepo := mock.NewMockRepository(ctrl)

	exe := func (usecase userUsecase, newUser *User) error {
		return usecase.Store(newUser)
	}

	t.Run("Test OK", func (t1 *testing.T) {
		newUser := User{Username: "newUser", Email: "n4@ya.ru", Password: "123456"}

		userRepo.EXPECT().GetByEmail(newUser.Email).Return(nil, vars.ErrNotFound)
		userRepo.EXPECT().Store(&newUser).Return(nil)
		usecase := configUserUsecase(userRepo)

		err := exe(usecase, &newUser)

		if err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error ErrAlreadyExist", func (t2 *testing.T) {
		newUser := User{Username: "newUser", Email: "m1@ya.ru", Password: "123456"}

		userRepo.EXPECT().GetByEmail("m1@ya.ru").Return(users.User[0], nil)
		usecase := configUserUsecase(userRepo)

		err := exe(usecase, &newUser)

		if err != vars.ErrAlreadyExist {
			t2.Fail()
		}
	})

	t.Run("Test with some error", func (t3 *testing.T) {
		newUser := User{Username: "newUser", Email: "n4@ya.ru", Password: "123456"}

		userRepo.EXPECT().GetByEmail("n4@ya.ru").Return(nil, vars.ErrNotFound)
		userRepo.EXPECT().Store(&newUser).Return(fmt.Errorf("some error"))
		usecase := configUserUsecase(userRepo)

		err := exe(usecase, &newUser)

		if err != vars.ErrInternalServerError {
			t3.Fail()
		}
	})
}
