package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/playlist"
	"database/sql"
)

type PlaylistRepository struct {
	db *sql.DB
}

func NewPlaylistRepository(db *sql.DB) playlist.Repository {
	return &PlaylistRepository{
		db: db,
	}
}

func (plR *PlaylistRepository) Store(playlist *models.Playlist) error {
	return plR.db.QueryRow("INSERT INTO playlists (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id, photo",
		playlist.Name,
		playlist.Description,
		playlist.OwnerID,
	).Scan(&playlist.ID, &playlist.Photo)
}
