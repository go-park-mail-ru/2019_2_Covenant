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

	rows, err := ar.db.Query("select Al.id, Al.artist_id, Al.name, Al.photo, Al.year, Ar.name, Ar.id " +
		"from albums Al join artists Ar on Al.artist_id = Ar.id where lower(Al.name) like '%' || $1 || '%' " +
		"OR lower(Ar.name) like '%' || $1 || '%' limit $2",
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

		if err := rows.Scan(&a.ID, &a.ArtistID, &a.Name, &a.Photo, &a.Year, &a.Artist, &a.ArtistID); err != nil {
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

	rows, err := ar.db.Query("SELECT Al.id, Al.artist_id, Al.name, Al.photo, Al.year, Ar.name, Ar.id " +
		"FROM albums Al JOIN artists Ar ON Al.artist_id = Ar.id ORDER BY Al.name LIMIT $1 OFFSET $2",
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
			&a.Artist,
			&a.ArtistID,
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

	if err := ar.db.QueryRow("SELECT Al.id, Al.artist_id, Al.name, Al.photo, Al.year, Ar.name, Ar.id " +
		"FROM albums Al JOIN artists Ar ON Al.artist_id = Ar.id WHERE Al.id = $1",
		id,
	).Scan(
		&a.ID,
		&a.ArtistID,
		&a.Name,
		&a.Photo,
		&a.Year,
		&a.Artist,
		&a.ArtistID,
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
	if err := ar.db.QueryRow("SELECT id FROM tracks WHERE album_id = $1 AND name = $2",
		albumID,
		track.Name,
	).Scan(); err == nil {
		return ErrAlreadyExist
	}

	if _, err := ar.db.Exec("INSERT INTO tracks (album_id, name, duration, path) VALUES ($1, $2, $3, $4)",
		track.AlbumID,
		track.Name,
		track.Duration,
		track.Path,
	); err != nil {
		return err
	}

	return nil
}

func (ar *AlbumRepository) GetTracksFrom(albumID uint64, authID uint64) ([]*models.Track, error) {
	var tracks []*models.Track

	rows, err := ar.db.Query(
		"select T.id, T.album_id, T.name, T.duration, T.path, Ar.name, Al.name, Ar.id, " +
			"T.id in (select track_id from favourites where user_id = $1) as favourite, " +
			"T.id in (select track_id from likes where user_id = $1) AS liked from tracks T " +
			"join albums Al ON T.album_id=Al.id " +
			"join artists Ar ON Al.artist_id=Ar.id where Al.id = $2;",
			authID, albumID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := &models.Track{}
		isFavourite := new(bool)
		isLiked := new(bool)

		if err := rows.Scan(&t.ID, &t.AlbumID, &t.Name, &t.Duration, &t.Path, &t.Artist,
				&t.Album, &t.ArtistID, isFavourite, isLiked); err != nil {
			return nil, err
		}

		if authID != 0 {
			t.IsFavourite = isFavourite
			t.IsLiked = isLiked
		}

		tracks = append(tracks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (ar *AlbumRepository) UpdatePhoto(albumID uint64, path string) error {
	if err := ar.db.QueryRow("UPDATE albums SET photo = $1 WHERE id = $2 RETURNING id",
		path,
		albumID,
	).Scan(&albumID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}
