package repository

import (
	"2019_2_Covenant/internal/artist"
	"2019_2_Covenant/internal/models"
	. "2019_2_Covenant/tools/vars"
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

func (ar *ArtistRepository) Fetch(count uint64, offset uint64) ([]*models.Artist, uint64, error) {
	var artists []*models.Artist
	var total uint64

	if err := ar.db.QueryRow("SELECT COUNT(*) FROM artists").Scan(&total); err != nil {
		return nil, total, err
	}

	rows, err := ar.db.Query("SELECT id, name, photo FROM artists ORDER BY name LIMIT $1 OFFSET $2",
		count,
		offset,
	)

	if err != nil {
		return nil, total, err
	}

	defer rows.Close()

	for rows.Next() {
		a := &models.Artist{}

		if err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.Photo,
		); err != nil {
			return nil, total, err
		}

		artists = append(artists, a)
	}

	if err := rows.Err(); err != nil {
		return nil, total, err
	}

	return artists, total, nil
}

func (ar *ArtistRepository) FindLike(name string, count uint64) ([]*models.Artist, error) {
	var artists []*models.Artist

	rows, err := ar.db.Query("select id, name, photo from artists where lower(name) like '%' || $1 || '%' limit $2",
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

		if err := rows.Scan(&a.ID, &a.Name, &a.Photo); err != nil {
			return nil, err
		}

		artists = append(artists, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}

func (ar *ArtistRepository) Store(artist *models.Artist) error {
	return ar.db.QueryRow("INSERT INTO artists (name) VALUES ($1) RETURNING id, photo",
		artist.Name,
	).Scan(&artist.ID, &artist.Photo)
}

func (ar *ArtistRepository) DeleteByID(id uint64) error {
	if err := ar.db.QueryRow("DELETE FROM artists WHERE id = $1 RETURNING id",
		id,
	).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (ar *ArtistRepository) UpdateByID(id uint64, name string) error {
	if err := ar.db.QueryRow("UPDATE artists SET name = $1 WHERE id = $2 RETURNING id",
		name,
		id,
	).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (ar *ArtistRepository) CreateAlbum(album *models.Album) error {
	return ar.db.QueryRow("INSERT INTO albums (artist_id, name, year) VALUES ($1, $2, $3) RETURNING id, photo",
		album.ArtistID, album.Name, album.Year,
	).Scan(&album.ID, &album.Photo)
}

func (ar *ArtistRepository) GetByID(id uint64) (*models.Artist, uint64, error) {
	a := &models.Artist{}
	var amountOfAlbums uint64

	if err := ar.db.QueryRow("SELECT id, name, photo FROM artists WHERE id = $1",
		id,
	).Scan(
		&a.ID,
		&a.Name,
		&a.Photo,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, amountOfAlbums, ErrNotFound
		}

		return nil, amountOfAlbums, err
	}

	if err := ar.db.QueryRow("SELECT COUNT(*) FROM albums WHERE artist_id = $1",
		id,
	).Scan(&amountOfAlbums); err != nil {
		if err == sql.ErrNoRows {
			return a, amountOfAlbums, nil
		}

		return nil, amountOfAlbums, err
	}

	return a, amountOfAlbums, nil
}

func (ar *ArtistRepository) UpdatePhoto(artistID uint64, path string) error {
	if err := ar.db.QueryRow("UPDATE artists SET photo = $1 WHERE id = $2 RETURNING id",
		path,
		artistID,
	).Scan(&artistID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (ar *ArtistRepository) GetArtistAlbums(artistID uint64, count uint64, offset uint64) ([]*models.Album, uint64, error) {
	var albums []*models.Album
	var total uint64

	if err := ar.db.QueryRow("SELECT COUNT(*) FROM albums WHERE artist_id = $1",
		artistID,
	).Scan(&total); err != nil {
		return nil, total, err
	}

	rows, err := ar.db.Query("SELECT id, name, photo, year FROM albums "+
		"WHERE artist_id = $1 ORDER BY name LIMIT $2 OFFSET $3",
		artistID,
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

func (ar *ArtistRepository) GetTracks(artistID uint64, count uint64, offset uint64, authID uint64) ([]*models.Track, uint64, error) {
	var tracks []*models.Track
	var total uint64

	if err := ar.db.QueryRow("SELECT COUNT(*) FROM tracks T JOIN albums Al ON T.album_id=Al.id "+
		"JOIN artists Ar ON Al.artist_id=Ar.id WHERE Ar.id = $1", artistID).Scan(&total); err != nil {
		return nil, total, err
	}

	rows, err := ar.db.Query(
		"SELECT T.id, T.album_id, T.name, T.duration, Al.photo, Al.name, T.path, "+
			"T.id in (select track_id from favourites where user_id = $1) as favourite, " +
			"T.id in (select track_id from likes where user_id = $1) AS liked FROM tracks T "+
			"JOIN albums Al ON T.album_id = Al.id "+
			"JOIN artists Ar ON Al.artist_id = Ar.id WHERE Ar.id = $2 LIMIT $3 OFFSET $4",
		authID, artistID, count, offset,
	)

	if err != nil {
		return nil, total, err
	}

	defer rows.Close()

	for rows.Next() {
		t := &models.Track{}
		isFavourite := new(bool)
		isLiked := new(bool)

		if err := rows.Scan(&t.ID, &t.AlbumID, &t.Name, &t.Duration, &t.Photo, &t.Album, &t.Path, isFavourite, isLiked); err != nil {
			return nil, total, err
		}

		if authID != 0 {
			t.IsFavourite = isFavourite
			t.IsLiked = isLiked
		}

		tracks = append(tracks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, total, err
	}

	return tracks, total, nil
}
