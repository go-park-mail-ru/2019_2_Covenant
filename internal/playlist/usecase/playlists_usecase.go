package usecase

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/playlist"
)

type PlaylistUsecase struct {
	playlistRepo playlist.Repository
}

func NewPlaylistUsecase(repo playlist.Repository) playlist.Usecase {
	return &PlaylistUsecase{
		playlistRepo: repo,
	}
}

func (pUC *PlaylistUsecase) Store(playlist *models.Playlist) error {
	err := pUC.playlistRepo.Store(playlist)

	if err != nil {
		return err
	}

	return nil
}
