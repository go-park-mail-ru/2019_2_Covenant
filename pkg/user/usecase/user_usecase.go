package usecase

import (
	files "2019_2_Covenant/pkg/file_processor"
	"2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/user/repository"
	. "2019_2_Covenant/tools/vars"
	"context"
	"io"
)

type UserUsecase struct {
	userRepo repository.UsersClient
	filesRepo files.Repository
}

func NewUserUsecase(ur repository.UsersClient, fr files.Repository) *UserUsecase {
	return &UserUsecase{
		userRepo: ur,
		filesRepo: fr,
	}
}

func toModel(user *repository.User) *models.User {
	return &models.User{
		ID:            user.Id,
		Nickname:      user.Nickname,
		Email:         user.Email,
		PlainPassword: user.PlainPassword,
		Password:      user.Password,
		Avatar:        user.Avatar,
		Role:          int8(user.Role),
		Access:        int8(user.Access),
		Subscription:  &user.Subscription,
	}
}

func fromModel(user *models.User) *repository.User {

	return &repository.User{
		Id:            user.ID,
		Nickname:      user.Nickname,
		Email:         user.Email,
		PlainPassword: user.PlainPassword,
		Password:      user.Password,
		Avatar:        user.Avatar,
		Role:          int32(user.Role),
		Access:        int32(user.Access),
		Subscription:  false,
	}
}

func (uUC *UserUsecase) Fetch(ctx context.Context, count uint64) ([]*models.User, error) {
	userArray, err := uUC.userRepo.Fetch(ctx, &repository.FetchRequest{Count:count})
	if err != nil {
		return nil, nil
	}

	users := make([]*models.User, 0, len(userArray.Users))
	for _, user := range userArray.Users {
		users = append(users, toModel(user))
	}

	return users, nil
}

func (uUC *UserUsecase) Store(ctx context.Context, newUser *models.User) error {
	if err := newUser.BeforeStore(); err != nil {
		return ErrInternalServerError
	}

	user, err := uUC.userRepo.Store(ctx, fromModel(newUser))
	if err != nil {
		return ErrAlreadyExist
	}

	*newUser = *toModel(user)
	return nil
}

func (uUC *UserUsecase) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	usr, err := uUC.userRepo.GetByEmail(ctx, &repository.GetByEmailRequest{Email:email})
	if err != nil {
		return nil, err
	}

	return toModel(usr), nil
}

func (uUC *UserUsecase) GetByID(ctx context.Context, userID uint64) (*models.User, error) {
	usr, err := uUC.userRepo.GetByID(ctx, &repository.GetByIDRequest{Id: userID})
	if err != nil {
		return nil, err
	}

	return toModel(usr), nil
}

func (uUC *UserUsecase) GetByNickname(ctx context.Context, nickname string, authID uint64) (*models.User, error) {
	usr, err := uUC.userRepo.GetByNickname(ctx, &repository.GetByNicknameRequest{Nickname:nickname, AuthID:authID})
	if err != nil {
		return nil, ErrNotFound
	}

	return toModel(usr), nil
}

func (uUC *UserUsecase) UpdateAvatar(ctx context.Context, userID uint64, photo io.Reader) error {
	return uUC.filesRepo.ProcessAvatar(ctx, photo, userID)
}

func (uUC *UserUsecase) UpdatePassword(ctx context.Context, id uint64, plainPassword string) error {
	password, err := models.EncryptPassword(plainPassword)

	if err != nil {
		return ErrInternalServerError
	}

	if _, err := uUC.userRepo.UpdatePassword(ctx, &repository.UpdatePasswordRequest{Id:id, Password:password} ); err != nil {
		return ErrInternalServerError
	}

	return nil
}

func (uUC *UserUsecase) Update(ctx context.Context, id uint64, nickname string, email string) (*models.User, error) {
	usr, err := uUC.userRepo.Update(ctx, &repository.UpdateRequest{Id:id, Nickname:nickname, Email:email})
	if err != nil {
		return nil, ErrAlreadyExist
	}

	return toModel(usr), nil
}

func (uUC *UserUsecase) GetFollowing(ctx context.Context, id uint64, count uint64, offset uint64) ([]*models.User, uint64, error) {
	response, err := uUC.userRepo.GetFollowing(ctx, &repository.GetFollowRequest{Id: id, Count:count, Offset:offset})

	if err != nil {
		return nil, 0, err
	}

	users := make([]*models.User, 0, len(response.Users))
	for _, user := range response.Users {
		users = append(users, toModel(user))
	}

	return users, response.Total, nil
}

func (uUC *UserUsecase) GetFollowers(ctx context.Context, id uint64, count uint64, offset uint64) ([]*models.User, uint64, error) {
	response, err := uUC.userRepo.GetFollowers(ctx, &repository.GetFollowRequest{Id: id, Count:count, Offset:offset})

	if err != nil {
		return nil, 0, err
	}

	users := make([]*models.User, 0, len(response.Users))
	for _, user := range response.Users {
		users = append(users, toModel(user))
	}

	return users, response.Total, nil
}

func (uUC *UserUsecase) FindLike(ctx context.Context, name string, count uint64) ([]*models.User, error) {
	userArray, err := uUC.userRepo.FindLike(ctx, &repository.FindLikeRequest{Name: name, Count: count})
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0, len(userArray.Users))
	for _, user := range userArray.Users {
		users = append(users, toModel(user))
	}

	return users, nil
}
