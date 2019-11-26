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

func (plR *PlaylistRepository) Fetch(userID uint64, count uint64, offset uint64) ([]*models.Playlist, uint64, error) {
	var playlists []*models.Playlist
	var total uint64

	if err := plR.db.QueryRow("SELECT COUNT(*) FROM playlists WHERE owner_id = $1",
		userID,
	).Scan(
		&total,
	); err != nil {
		return nil, total, err
	}

	rows, err := plR.db.Query("SELECT id, name, description, photo FROM playlists WHERE owner_id = $1 LIMIT $2 OFFSET $3",
		userID,
		count,
		offset,
	)

	if err != nil {
		return nil, total, err
	}

	defer rows.Close()

	for rows.Next() {
		p := &models.Playlist{}

		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Photo,
		); err != nil {
			return nil, total, err
		}

		p.OwnerID = userID

		playlists = append(playlists, p)
	}

	if err := rows.Err(); err != nil {
		return nil, total, err
	}

	return playlists, total, nil
}
