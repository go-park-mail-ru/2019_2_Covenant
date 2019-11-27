package repository

import (
	"2019_2_Covenant/internal/artist"
	"2019_2_Covenant/internal/models"
	"database/sql"
	"strings"
)

type ArtistRepository struct {
	db *sql.DB
}

func NewArtistRepository(db *sql.DB) artist.Repository {
	return &ArtistRepository{
		db: db,
	}
}

func (ar *ArtistRepository) FindLike(name string, count uint64) ([]*models.Artist, error) {
	var artists []*models.Artist

	rows, err := ar.db.Query("select id, name from artists where lower(name) like '%' || $1 || '%' limit $2",
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
		a := &models.Artist{}

		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}

		artists = append(artists, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}
