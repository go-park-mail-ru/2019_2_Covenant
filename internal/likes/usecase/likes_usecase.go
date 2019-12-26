package usecase

import (
	"2019_2_Covenant/internal/likes"
	. "2019_2_Covenant/tools/vars"
)

type LikesUsecase struct {
	likesRepo likes.Repository
}

func NewLikesUsecase(repo likes.Repository) likes.Usecase {
	return &LikesUsecase{
		likesRepo: repo,
	}
}

func (lUC *LikesUsecase) Like(userID uint64, trackID uint64) error {
	err := lUC.likesRepo.Like(userID, trackID)

	if err == ErrAlreadyExist {
		return err
	}

	if err != nil {
		return ErrInternalServerError
	}

	return nil
}

func (lUC *LikesUsecase) Unlike(userID uint64, trackID uint64) error {
	err := lUC.likesRepo.Unlike(userID, trackID)

	if err == ErrNotFound {
		return err
	}

	if err != nil {
		return ErrInternalServerError
	}

	return nil
}
