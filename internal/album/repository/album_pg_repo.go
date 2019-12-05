package repository

import (
	"2019_2_Covenant/internal/album"
	"2019_2_Covenant/internal/models"
	. "2019_2_Covenant/tools/vars"
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

func (ar *AlbumRepository) DeleteByID(id uint64) error {
	if err := ar.db.QueryRow("DELETE FROM albums WHERE id = $1 RETURNING id",
		id,
	).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (ar *AlbumRepository) UpdateByID(albumID uint64, artistID uint64, name string, year string) error {
	if err := ar.db.QueryRow("UPDATE albums SET artist_id = $1, name = $2, year = $3 WHERE id = $4 RETURNING id",
		artistID,
		name,
		year,
		albumID,
	).Scan(&albumID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (ar *AlbumRepository) Fetch(count uint64, offset uint64) ([]*models.Album, uint64, error) {
	var albums []*models.Album
	var total uint64

	if err := ar.db.QueryRow("SELECT COUNT(*) FROM albums").Scan(&total); err != nil {
		return nil, total, err
	}

	rows, err := ar.db.Query("SELECT id, artist_id, name, photo, year FROM albums ORDER BY name LIMIT $1 OFFSET $2",
		count,
		offset,
	)

	if err != nil {
		return nil, total, err
	}

	defer rows.Close()

	for rows.Next() {
		a := &models.Album{}

		if err := rows.Scan(
			&a.ID,
			&a.ArtistID,
			&a.Name,
			&a.Photo,
			&a.Year,
		); err != nil {
			return nil, total, err
		}

		albums = append(albums, a)
	}

	if err := rows.Err(); err != nil {
		return nil, total, err
	}

	return albums, total, nil
}

func (ar *AlbumRepository) GetByID(id uint64) (*models.Album, uint64, error) {
	a := &models.Album{}
	var amountOfTracks uint64

	if err := ar.db.QueryRow("SELECT id, artist_id, name, photo, year FROM albums WHERE id = $1",
		id,
	).Scan(
		&a.ID,
		&a.ArtistID,
		&a.Name,
		&a.Photo,
		&a.Year,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, amountOfTracks, ErrNotFound
		}

		return nil, amountOfTracks, err
	}

	if err := ar.db.QueryRow("SELECT COUNT(*) FROM tracks WHERE album_id = $1",
		id,
	).Scan(&amountOfTracks); err != nil {
		if err == sql.ErrNoRows {
			return a, amountOfTracks, nil
		}

		return nil, amountOfTracks, err
	}

	return a, amountOfTracks, nil
}

func (ar *AlbumRepository) AddTrack(albumID uint64, track *models.Track) error {
	if err := ar.db.QueryRow("SELECT id FROM tracks WHERE album_id = $1",
		albumID,
	).Scan(); err == nil {
		return ErrAlreadyExist
	}

	if _, err := ar.db.Exec("INSERT INTO tracks (album_id, name, duration) VALUES ($1, $2, $3)",
		track.AlbumID,
		track.Name,
		track.Duration,
	); err != nil {
		return err
	}

	return nil
}
