package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/track"
	"2019_2_Covenant/internal/vars"
	"database/sql"
)

// Для тестирования только этого файла:
// go test -v -cover -race ./internal/track/repository

type TrackRepository struct {
	db *sql.DB
}

func NewTrackRepository(db *sql.DB) track.Repository {
	return &TrackRepository{
		db: db,
	}
}

func (tr *TrackRepository) Fetch(count uint64, offset uint64) ([]*models.Track, error) {
	var tracks []*models.Track

	rows, err := tr.db.Query(
		"SELECT T.id, T.album_id, Ar.id, T.name, T.duration, Al.photo, Ar.name, Al.name, T.path FROM tracks T " +
		"JOIN albums Al ON T.album_id = Al.id " +
		"JOIN artists Ar ON Al.artist_id = Ar.id LIMIT $1 OFFSET $2",
		count,
		offset)

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
			ID:       t.ID,
			AlbumID:  t.AlbumID,
			ArtistID: t.ArtistID,
			Name:     t.Name,
			Duration: t.Duration,
			Photo:    t.Photo,
			Artist:   t.Artist,
			Album:    t.Album,
			Path:     t.Path,
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

func (th *TrackRepository) FetchFavourites(userID uint64, count uint64, offset uint64) ([]*models.Track, uint64, error) {
	var tracks []*models.Track
	var total uint64

	if err := th.db.QueryRow("SELECT COUNT(*) FROM tracks T JOIN favourites F on T.id = F.track_id WHERE F.user_id = $1",
		userID,
	).Scan(
		&total,
	); err != nil {
		return nil, total, err
	}

	rows, err := th.db.Query(
		"SELECT T.id, T.album_id, Ar.id, T.name, T.duration, Al.photo, Ar.name, Al.name, T.path FROM tracks T "+
			"JOIN favourites F ON T.id = F.track_id "+
			"JOIN albums Al ON T.album_id = Al.id "+
			"JOIN artists Ar ON Al.artist_id = Ar.id "+
			"WHERE F.user_id = $1 LIMIT $2",
		userID,
		count,
	)

	if err != nil {
		return nil, total, err
	}

	defer rows.Close()

	for rows.Next() {
		t := &models.Track{}

		if err := rows.Scan(&t.ID, &t.AlbumID, &t.ArtistID, &t.Name, &t.Duration,
			&t.Photo, &t.Artist, &t.Album, &t.Path,
		); err != nil {
			return nil, total, err
		}

		tracks = append(tracks, &models.Track{
			ID:       t.ID,
			AlbumID:  t.AlbumID,
			ArtistID: t.ArtistID,
			Name:     t.Name,
			Duration: t.Duration,
			Photo:    t.Photo,
			Artist:   t.Artist,
			Album:    t.Album,
			Path:     t.Path,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, total, err
	}

	return tracks, total, nil
}
