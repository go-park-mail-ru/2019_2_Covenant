package usecase

import (
	. "2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/user"
	mock "2019_2_Covenant/internal/user/mocks"
	"2019_2_Covenant/tools/vars"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
)

//go:generate mockgen -source=../repository.go -destination=../mocks/mock_repository.go -package=mock

type Users struct {
	User []*User
}

var users = Users{
	User: []*User{
		{ID: 1, Nickname: "marshal", Email: "m1@ya.ru", Password: "123456"},
		{ID: 2, Nickname: "plaksenka", Email: "p2@ya.ru", Password: "123456"},
		{ID: 3, Nickname: "svya", Email: "s3@ya.ru", Password: "123456"},
	},
}

func configUserUsecase(userRepo *mock.MockRepository) user.Usecase {
	return NewUserUsecase(userRepo)
}

func TestUserUsecase_Fetch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase user.Usecase, count uint64) ([]*User, error) {
		return usecase.Fetch(count)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		count := uint64(3)

		userRepo.EXPECT().Fetch(count).Return(users.User, nil)
		usecase := configUserUsecase(userRepo)

		expUsers, err := exe(usecase, count)

		if gomock.Not(users.User).Matches(expUsers) || err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		count := uint64(3)

		userRepo.EXPECT().Fetch(count).Return(nil, fmt.Errorf("internal error"))
		usecase := configUserUsecase(userRepo)

		expUsers, err := exe(usecase, count)

		if expUsers != nil || err == nil {
			t2.Fail()
		}
	})
}

func TestUserUsecase_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase user.Usecase, email string) (*User, error) {
		return usecase.GetByEmail(email)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		setEmail := "m1@ya.ru"

		userRepo.EXPECT().GetByEmail(setEmail).Return(users.User[0], nil)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, setEmail)

		if expUser != users.User[0] || err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		setEmail := "notfound@ya.ru"

		userRepo.EXPECT().GetByEmail(setEmail).Return(nil, vars.ErrNotFound)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, setEmail)

		if expUser != nil || err != vars.ErrNotFound {
			t2.Fail()
		}
	})
}

func TestUserUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase user.Usecase, ID uint64) (*User, error) {
		return usecase.GetByID(ID)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		setID := uint64(1)

		userRepo.EXPECT().GetByID(setID).Return(users.User[0], nil)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, setID)

		if expUser != users.User[0] || err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		setID := uint64(5)
		userRepo.EXPECT().GetByID(setID).Return(nil, vars.ErrNotFound)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, setID)

		if expUser != nil || err != vars.ErrNotFound {
			t2.Fail()
		}
	})
}

func TestUserUsecase_GetByNickname(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase user.Usecase, nickname string) (*User, error) {
		return usecase.GetByNickname(nickname)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		setNickname := "newMarshal"

		userRepo.EXPECT().GetByNickname(setNickname).Return(users.User[0], nil)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, setNickname)

		if expUser != users.User[0] || err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		setEmail := "notfound@ya.ru"

		userRepo.EXPECT().GetByNickname(setEmail).Return(nil, vars.ErrNotFound)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, setEmail)

		if expUser != nil || err != vars.ErrNotFound {
			t2.Fail()
		}
	})
}

func TestUserUsecase_Store(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase user.Usecase, newUser *User) error {
		return usecase.Store(newUser)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		newUser := User{Nickname: "newUser", Email: "n4@ya.ru", Password: "123456"}

		userRepo.EXPECT().GetByEmail(newUser.Email).Return(nil, vars.ErrNotFound)
		userRepo.EXPECT().Store(&newUser).Return(nil)
		usecase := configUserUsecase(userRepo)

		err := exe(usecase, &newUser)

		if err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error ErrAlreadyExist", func(t2 *testing.T) {
		newUser := User{Nickname: "newUser", Email: "m1@ya.ru", Password: "123456"}

		userRepo.EXPECT().GetByEmail("m1@ya.ru").Return(users.User[0], nil)
		usecase := configUserUsecase(userRepo)

		err := exe(usecase, &newUser)

		if err != vars.ErrAlreadyExist {
			t2.Fail()
		}
	})

	t.Run("Test with some error", func(t3 *testing.T) {
		newUser := User{Nickname: "newUser", Email: "n4@ya.ru", Password: "123456"}

		userRepo.EXPECT().GetByEmail("n4@ya.ru").Return(nil, vars.ErrNotFound)
		userRepo.EXPECT().Store(&newUser).Return(fmt.Errorf("some error"))
		usecase := configUserUsecase(userRepo)

		err := exe(usecase, &newUser)

		if err != vars.ErrInternalServerError {
			t3.Fail()
		}
	})
}

func TestUserUsecase_UpdateAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase user.Usecase, ID uint64, avatarPath string) (*User, error) {
		return usecase.UpdateAvatar(ID, avatarPath)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		ID := uint64(1)
		avatarPath := "some path"

		userRepo.EXPECT().UpdateAvatar(ID, avatarPath).Return(users.User[0], nil)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, ID, avatarPath)

		if expUser != users.User[0] || err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		ID := uint64(5)
		avatarPath := "some path"

		userRepo.EXPECT().UpdateAvatar(ID, avatarPath).Return(nil, fmt.Errorf("some error"))
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, ID, avatarPath)

		if expUser != nil || err != vars.ErrInternalServerError {
			t2.Fail()
		}
	})
}

func TestUserUsecase_UpdateNickname(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase user.Usecase, ID uint64, nickname string, email string) (*User, error) {
		return usecase.Update(ID, nickname, email)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		ID := uint64(1)
		nickname := "some nickname"
		email := "some email"

		userRepo.EXPECT().Update(ID, nickname, email).Return(users.User[0], nil)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, ID, nickname, email)

		if expUser != users.User[0] || err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		ID := uint64(5)
		nickname := "some nickname"
		email := "some email"

		userRepo.EXPECT().Update(ID, nickname, email).Return(nil, fmt.Errorf("some error"))
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, ID, nickname, email)

		if expUser != nil || err == nil {
			t2.Fail()
		}
	})

	t.Run("Test with ErrAlreadyExist", func(t3 *testing.T) {
		ID := uint64(5)
		nickname := "some nickname"
		email := "some email"

		userRepo.EXPECT().Update(ID, nickname, email).Return(nil, vars.ErrAlreadyExist)
		usecase := configUserUsecase(userRepo)

		expUser, err := exe(usecase, ID, nickname, email)

		if expUser != nil || err != vars.ErrAlreadyExist {
			t3.Fail()
		}
	})
}

func TestUserUsecase_UpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase user.Usecase, ID uint64, password string) error {
		return usecase.UpdatePassword(ID, password)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		ID := uint64(1)
		password := "some email"

		userRepo.EXPECT().UpdatePassword(ID, gomock.Any()).Return(nil)
		usecase := configUserUsecase(userRepo)

		err := exe(usecase, ID, password)

		if err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		ID := uint64(5)
		password := ""

		userRepo.EXPECT().UpdatePassword(ID, gomock.Any()).Return(fmt.Errorf("some error"))
		usecase := configUserUsecase(userRepo)

		err := exe(usecase, ID, password)

		if err == nil {
			t2.Fail()
		}
	})
}
