package usecase

import (
	. "2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/user"
	"2019_2_Covenant/pkg/vars"
)

type userUsecase struct {
	userRepo user.Repository
}

func NewUserUsecase(ur user.Repository) user.Usecase {
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
		return vars.ErrAlreadyExist
	}

	err := uUC.userRepo.Store(newUser)

	if err != nil {
		return vars.ErrInternalServerError
	}

	return nil
}

func (uUC *userUsecase) GetByEmail(email string) (*User, error) {
	usr, err := uUC.userRepo.GetByEmail(email)

	if err != nil {
		return nil, vars.ErrNotFound
	}

	return usr, nil
}

func (uUC *userUsecase) GetByID(userID uint64) (*User, error) {
	usr, err := uUC.userRepo.GetByID(userID)

	if err != nil {
		return nil, vars.ErrNotFound
	}

	return usr, nil
}
