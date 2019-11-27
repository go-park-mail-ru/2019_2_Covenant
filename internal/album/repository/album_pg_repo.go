package repository

import (
	"2019_2_Covenant/internal/album"
	"2019_2_Covenant/internal/models"
	"database/sql"
	"strings"
)

type AlbumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) album.Repository {
	return &AlbumRepository{
		db: db,
	}
}

func (ar *AlbumRepository) FindLike(name string, count uint64) ([]*models.Album, error) {
	var albums []*models.Album

	rows, err := ar.db.Query("select id, artist_id, name, photo, year from albums where lower(name) like '%' || $1 || '%' limit $2",
		strings.ToLower(name),
		count)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		a := &models.Album{}

		if err := rows.Scan(&a.ID, &a.ArtistID, &a.Name, &a.Photo, &a.Year); err != nil {
			return nil, err
		}

		albums = append(albums, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}
