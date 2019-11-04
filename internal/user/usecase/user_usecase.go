package usecase

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/user"
	"2019_2_Covenant/internal/vars"
)

type userUsecase struct {
	userRepo user.Repository
}

func NewUserUsecase(ur user.Repository) user.Usecase {
	return &userUsecase{
		userRepo: ur,
	}
}

func (uUC *userUsecase) Fetch(count uint64) ([]*models.User, error) {
	users, err := uUC.userRepo.Fetch(count)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uUC *userUsecase) Store(newUser *models.User) (*models.User, error) {
	exist, _ := uUC.userRepo.GetByEmail(newUser.Email)

	if exist != nil {
		return nil, vars.ErrAlreadyExist
	}

	usr, err := uUC.userRepo.Store(newUser)

	if err != nil {
		return nil, vars.ErrInternalServerError
	}

	return usr, nil
}

func (uUC *userUsecase) GetByEmail(email string) (*models.User, error) {
	usr, err := uUC.userRepo.GetByEmail(email)

	if err != nil {
		return nil, vars.ErrNotFound
	}

	return usr, nil
}

func (uUC *userUsecase) GetByID(userID uint64) (*models.User, error) {
	usr, err := uUC.userRepo.GetByID(userID)

	if err != nil {
		return nil, vars.ErrNotFound
	}

	return usr, nil
}

func (uUC *userUsecase) UpdateAvatar(id uint64, avatarPath string) (*models.User, error) {
	usr, err := uUC.userRepo.UpdateAvatar(id, avatarPath)

	if err != nil {
		return nil, vars.ErrInternalServerError
	}

	return usr, nil
}

func (uUC *userUsecase) UpdateNickname(id uint64, nickname string) (*models.User, error) {
	usr, err := uUC.userRepo.UpdateNickname(id, nickname)

	if err == vars.ErrAlreadyExist {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return usr, nil
}
