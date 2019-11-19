package usecase

import (
	. "2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/session"
	mock "2019_2_Covenant/internal/session/mocks"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

//go:generate mockgen -source=../repository.go -destination=../mocks/mock_repository.go -package=mock
//go:generate mockgen -source=../../models/user.go -destination=../../models/mocks/mock_user.go -package=mock

type Sessions struct {
	Session []*Session
}

var sessions = Sessions{
	Session: []*Session{
		{ID: 1, UserID: 1, Expires: time.Now().Add(24 * time.Hour), Data: uuid.New().String()},
		{ID: 2, UserID: 2, Expires: time.Now().Add(1 * time.Hour), Data: uuid.New().String()},
		{ID: 3, UserID: 3, Expires: time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC), Data: uuid.New().String()},
	},
}

func configSessionUsecase(sessRepo *mock.MockRepository) session.Usecase {
	return NewSessionUsecase(sessRepo)
}

func TestSessionUsecase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase session.Usecase, value string) (*Session, error) {
		return usecase.Get(value)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		sessRepo.EXPECT().Get(sessions.Session[0].Data).Return(sessions.Session[0], nil)
		usecase := configSessionUsecase(sessRepo)

		expSess, err := exe(usecase, sessions.Session[0].Data)

		if expSess != sessions.Session[0] && err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with Error", func(t2 *testing.T) {
		sessValueNotFound := uuid.New().String()
		sessRepo.EXPECT().Get(sessValueNotFound).Return(nil, fmt.Errorf("some error"))
		usecase := configSessionUsecase(sessRepo)

		expSess, err := exe(usecase, sessValueNotFound)

		if expSess != nil && err == nil {
			t2.Fail()
		}
	})
}

func TestSessionUsecase_DeleteByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase session.Usecase, id uint64) error {
		return usecase.DeleteByID(id)
	}

	t.Run("Test OK", func(t1 *testing.T) {
		sessRepo.EXPECT().DeleteByID(sessions.Session[2].ID).Return(nil)
		usecase := configSessionUsecase(sessRepo)

		err := exe(usecase, sessions.Session[2].ID)

		if err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		Id := uint64(1)
		sessRepo.EXPECT().DeleteByID(Id).Return(fmt.Errorf("some error"))
		usecase := configSessionUsecase(sessRepo)

		err := exe(usecase, Id)

		if err == nil {
			t2.Fail()
		}
	})
}

func TestSessionUsecase_Store(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessRepo := mock.NewMockRepository(ctrl)

	exe := func(usecase session.Usecase, newSess *Session) error {
		return usecase.Store(newSess)
	}

	newSess := Session{UserID: 1, Expires: time.Now().Add(24 * time.Hour), Data: uuid.New().String()}
	t.Run("Test OK", func(t1 *testing.T) {
		sessRepo.EXPECT().Store(&newSess).Return(nil)
		usecase := configSessionUsecase(sessRepo)

		err := exe(usecase, &newSess)

		if err != nil {
			t1.Fail()
		}
	})

	t.Run("Test with error", func(t2 *testing.T) {
		sessRepo.EXPECT().Store(&newSess).Return(fmt.Errorf("some error"))
		usecase := configSessionUsecase(sessRepo)

		err := exe(usecase, &newSess)

		if err == nil {
			t2.Fail()
		}
	})
}
