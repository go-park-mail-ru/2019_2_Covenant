package usecase

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/user"
	. "2019_2_Covenant/tools/vars"
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

	if users == nil {
		users = []*models.User{}
	}

	return users, nil
}

func (uUC *userUsecase) Store(newUser *models.User) error {
	if err := newUser.BeforeStore(); err != nil {
		return ErrInternalServerError
	}

	if err := uUC.userRepo.Store(newUser); err != nil {
		return ErrAlreadyExist
	}

	return nil
}

func (uUC *userUsecase) GetByEmail(email string) (*models.User, error) {
	usr, err := uUC.userRepo.GetByEmail(email)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (uUC *userUsecase) GetByID(userID uint64) (*models.User, error) {
	usr, err := uUC.userRepo.GetByID(userID)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (uUC *userUsecase) GetByNickname(nickname string, authID uint64) (*models.User, error) {
	usr, err := uUC.userRepo.GetByNickname(nickname, authID)

	if err != nil {
		return nil, ErrNotFound
	}

	return usr, nil
}

func (uUC *userUsecase) UpdateAvatar(id uint64, avatarPath string) (*models.User, error) {
	usr, err := uUC.userRepo.UpdateAvatar(id, avatarPath)

	if err != nil {
		return nil, ErrInternalServerError
	}

	return usr, nil
}

func (uUC *userUsecase) UpdatePassword(id uint64, plainPassword string) error {
	password, err := models.EncryptPassword(plainPassword)

	if err != nil {
		return ErrInternalServerError
	}

	if err := uUC.userRepo.UpdatePassword(id, password); err != nil {
		return ErrInternalServerError
	}

	return nil
}

func (uUC *userUsecase) Update(id uint64, nickname string, email string) (*models.User, error) {
	usr, err := uUC.userRepo.Update(id, nickname, email)

	if err != nil {
		return nil, ErrAlreadyExist
	}

	return usr, nil
}

func (uUC *userUsecase) GetFollowing(id uint64, count uint64, offset uint64) ([]*models.User, uint64, error) {
	following, total, err := uUC.userRepo.GetFollowing(id, count, offset)

	if err != nil {
		return nil, total, err
	}

	if following == nil {
		following = []*models.User{}
	}

	return following, total, nil
}

func (uUC *userUsecase) GetFollowers(id uint64, count uint64, offset uint64) ([]*models.User, uint64, error) {
	followers, total, err := uUC.userRepo.GetFollowers(id, count, offset)

	if err != nil {
		return nil, total, err
	}

	if followers == nil {
		followers = []*models.User{}
	}

	return followers, total, nil
}

func (uUC *userUsecase) FindLike(name string, count uint64) ([]*models.User, error) {
	users, err := uUC.userRepo.FindLike(name, count)

	if err != nil {
		return nil, err
	}

	if users == nil {
		users = []*models.User{}
	}

	return users, nil
}
