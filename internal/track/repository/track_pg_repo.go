package repository

import (
	"2019_2_Covenant/internal/models"
	"2019_2_Covenant/internal/track"
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
		"SELECT T.id, T.album_id, Ar.id, T.name, T.duration, Al.photo, Ar.name, Al.name FROM tracks T " +
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
		)

		if err := rows.Scan(&id, &albumID, &artistID, &name, &duration, &photo, &artist, &album); err != nil {
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
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}
