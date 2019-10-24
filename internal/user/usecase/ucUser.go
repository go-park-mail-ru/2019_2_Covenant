package usecase

import (
	. "2019_2_Covenant/internal/models"
	user2 "2019_2_Covenant/internal/user"
	vars2 "2019_2_Covenant/internal/vars"
)

type userUsecase struct {
	userRepo user2.Repository
}

func NewUserUsecase(ur user2.Repository) user2.Usecase {
	return &userUsecase{
		userRepo: ur,
	}
}

func (uUC *userUsecase) FetchAll() ([]*User, error) {
	users, err := uUC.userRepo.FetchAll()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uUC *userUsecase) Store(newUser *User) error {
	exist, _ := uUC.userRepo.GetByEmail(newUser.Email)

	if exist != nil {
		return vars2.ErrAlreadyExist
	}

	err := uUC.userRepo.Store(newUser)

	if err != nil {
		return vars2.ErrInternalServerError
	}

	return nil
}

func (uUC *userUsecase) GetByEmail(email string) (*User, error) {
	usr, err := uUC.userRepo.GetByEmail(email)

	if err != nil {
		return nil, vars2.ErrNotFound
	}

	return usr, nil
}

func (uUC *userUsecase) GetByID(userID uint64) (*User, error) {
	usr, err := uUC.userRepo.GetByID(userID)

	if err != nil {
		return nil, vars2.ErrNotFound
	}

	return usr, nil
}
