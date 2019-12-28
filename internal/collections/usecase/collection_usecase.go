package usecase

import (
	"2019_2_Covenant/internal/collections"
	"2019_2_Covenant/internal/models"
	. "2019_2_Covenant/tools/vars"
)

type CollectionUsecase struct {
	collectionRepo collections.Repository
}

func NewCollectionUsecase(repo collections.Repository) collections.Usecase {
	return &CollectionUsecase{
		collectionRepo: repo,
	}
}

func (cu *CollectionUsecase) Store(collection *models.Collection) error {
	if err := cu.collectionRepo.Insert(collection); err != nil {
		return err
	}

	return nil
}

func (cu *CollectionUsecase) DeleteByID(id uint64) error {
	if err := cu.collectionRepo.DeleteByID(id); err != nil {
		return err
	}

	return nil
}

func (cu *CollectionUsecase) UpdateByID(collectionID uint64, name string, description string) error {
	if err := cu.collectionRepo.UpdateByID(collectionID, name, description); err != nil {
		return err
	}

	return nil
}

func (cu *CollectionUsecase) Fetch(count uint64, offset uint64) ([]*models.Collection, uint64, error) {
	albums, total, err := cu.collectionRepo.Select(count, offset)

	if err != nil {
		return nil, total, err
	}

	if albums == nil {
		albums = []*models.Collection{}
	}

	return albums, total, nil
}

func (cu *CollectionUsecase) GetByID(id uint64) (*models.Collection, uint64, error) {
	p, amountOfTracks, err := cu.collectionRepo.SelectByID(id)

	if err != nil {
		return nil, amountOfTracks, err
	}

	return p, amountOfTracks, nil
}

func (cu *CollectionUsecase) AddTrack(collectionID uint64, trackID uint64) error {
	err := cu.collectionRepo.InsertTrack(collectionID, trackID)

	if err == ErrAlreadyExist || err == ErrNotFound {
		return err
	}

	if err != nil {
		return ErrInternalServerError
	}

	return nil
}

func (cu *CollectionUsecase) GetTracks(collectionID uint64, authID uint64) ([]*models.Track, error) {
	tracks, err := cu.collectionRepo.SelectTracks(collectionID, authID)

	if err != nil {
		return nil, err
	}

	if tracks == nil {
		tracks = []*models.Track{}
	}

	return tracks, nil
}

func (cu *CollectionUsecase) UpdatePhoto(collectionID uint64, path string) error {
	if err := cu.collectionRepo.UpdatePhoto(collectionID, path); err != nil {
		return err
	}

	return nil
}
