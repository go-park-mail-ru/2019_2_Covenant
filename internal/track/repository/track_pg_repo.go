package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/track"
	"2019_2_Covenant/internal/vars"
	"database/sql"
)

type TrackRepository struct {
	db *sql.DB
}

func NewTrackRepository(db *sql.DB) track.Repository {
	return &TrackRepository{
		db: db,
	}
}

func (tr *TrackRepository) Fetch(count uint64) ([]*models.Track, error) {
	var tracks []*models.Track

	rows, err := tr.db.Query(
		"SELECT T.id, T.album_id, Ar.id, T.name, T.duration, Al.photo, Ar.name, Al.name, T.path FROM tracks T " +
		"JOIN albums Al ON T.album_id = Al.id " +
		"JOIN artists Ar ON Al.artist_id = Ar.id LIMIT $1",
		count)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id       uint64
			albumID  uint64
			artistID uint64
			name     string
			duration string
			photo    string
			artist   string
			album    string
			path     string
		)

		if err := rows.Scan(&id, &albumID, &artistID, &name, &duration, &photo, &artist, &album, &path); err != nil {
			return nil, err
		}

		tracks = append(tracks, &models.Track{
			ID: id,
			AlbumID: albumID,
			ArtistID: artistID,
			Name: name,
			Duration: duration,
			Photo: photo,
			Artist: artist,
			Album: album,
			Path: path,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (th *TrackRepository) StoreFavourite(userID uint64, trackID uint64) error {
	var favID uint64

	if err := th.db.QueryRow("SELECT id FROM favourites WHERE user_id = $1 AND track_id = $2",
		userID,
		trackID,
	).Scan(&favID); err == nil {
		return vars.ErrAlreadyExist
	}

	if _, err := th.db.Exec("INSERT INTO favourites (user_id, track_id) VALUES ($1, $2)",
		userID,
		trackID,
	); err != nil {
		return err
	}

	return nil
}

func (th *TrackRepository) RemoveFavourite(userID uint64, trackID uint64) error {
	var favID uint64

	if err := th.db.QueryRow("SELECT id FROM favourites WHERE user_id = $1 AND track_id = $2",
		userID,
		trackID,
	).Scan(&favID); err != nil {
		return vars.ErrNotFound
	}

	if _, err := th.db.Exec("DELETE FROM favourites WHERE id = $1",
		favID,
	); err != nil {
		return err
	}

	return nil
}

func (th *TrackRepository) FetchFavourites(userID uint64, count uint64) ([]*models.Track, error) {
	var tracks []*models.Track

	rows, err := th.db.Query(
		"SELECT T.id, T.album_id, Ar.id, T.name, T.duration, Al.photo, Ar.name, Al.name, T.path FROM tracks T " +
		"JOIN favourites F ON T.id = F.track_id " +
		"JOIN albums Al ON T.album_id = Al.id " +
		"JOIN artists Ar ON Al.artist_id = Ar.id " +
		"WHERE F.user_id = $1 LIMIT $2",
		userID,
		count,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := &models.Track{}

		if err := rows.Scan(&t.ID, &t.AlbumID, &t.ArtistID, &t.Name, &t.Duration,
			&t.Photo, &t.Artist, &t.Album, &t.Path,
		); err != nil {
			return nil, err
		}

		tracks = append(tracks, &models.Track{
			ID: t.ID,
			AlbumID: t.AlbumID,
			ArtistID: t.ArtistID,
			Name: t.Name,
			Duration: t.Duration,
			Photo: t.Photo,
			Artist: t.Artist,
			Album: t.Album,
			Path: t.Path,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}
