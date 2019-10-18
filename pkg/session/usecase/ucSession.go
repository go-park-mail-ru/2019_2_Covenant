package usecase

import (
	"2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/session"
)

type sessionUsecase struct {
	sessionRepo session.Repository
}

func NewSessionUsecase(sr session.Repository) session.Usecase {
	return &sessionUsecase{
		sessionRepo: sr,
	}
}

func (sUC sessionUsecase) Get(value string) (*models.Session, error) {
	sess, err := sUC.sessionRepo.Get(value)

	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (sUC *sessionUsecase) Store(newSession *models.Session) error {
	err := sUC.sessionRepo.Store(newSession)

	if err != nil {
		return err
	}

	return nil
}

func (sUC *sessionUsecase) DeleteByID(id uint64) error {
	err := sUC.sessionRepo.DeleteByID(id)

	if err != nil {
		return err
	}

	return nil
}
