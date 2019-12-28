package repository

import (
	"2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/playlist"
	. "2019_2_Covenant/tools/vars"
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

func (plR *PlaylistRepository) DeleteByID(playlistID uint64) error {
	if err := plR.db.QueryRow("DELETE FROM playlists WHERE id = $1 RETURNING id",
		playlistID,
	).Scan(&playlistID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (plR *PlaylistRepository) AddToPlaylist(playlistID uint64, trackID uint64) error {
	if err := plR.db.QueryRow("SELECT id FROM playlist_track WHERE playlist_id = $1 AND track_id = $2",
		playlistID,
		trackID,
	).Scan(); err == nil {
		return ErrAlreadyExist
	}

	if _, err := plR.db.Exec("INSERT INTO playlist_track (playlist_id, track_id) VALUES ($1, $2)",
		playlistID,
		trackID,
	); err != nil {
		return err
	}

	return nil
}

func (plR *PlaylistRepository) RemoveFromPlaylist(playlistID uint64, trackID uint64) error {
	res, err := plR.db.Exec("DELETE from playlist_track WHERE playlist_id = $1 AND track_id = $2",
		playlistID,
		trackID,
	)

	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (plR *PlaylistRepository) GetSinglePlaylist(playlistID uint64) (*models.Playlist, uint64, error) {
	p := &models.Playlist{}
	var amountOfTracks uint64

	if err := plR.db.QueryRow("SELECT id, name, description, photo, owner_id FROM playlists WHERE id = $1",
		playlistID,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Photo,
		&p.OwnerID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, amountOfTracks, ErrNotFound
		}

		return nil, amountOfTracks, err
	}

	if err := plR.db.QueryRow("SELECT COUNT(*) FROM playlist_track WHERE playlist_id = $1",
		playlistID,
	).Scan(&amountOfTracks); err != nil {
		if err == sql.ErrNoRows {
			return p, amountOfTracks, nil
		}

		return nil, amountOfTracks, err
	}

	return p, amountOfTracks, nil
}

func (plR *PlaylistRepository) GetTracksFrom(playlistID uint64, authID uint64) ([]*models.Track, error) {
	var tracks []*models.Track

	rows, err := plR.db.Query(
		"select T.id, T.name, T.duration, T.path, Ar.name, Ar.id, " +
			"T.id in (select track_id from favourites where user_id = $1) AS favourite, " +
			"T.id in (select track_id from likes where user_id = $1) AS liked from playlist_track PT " +
			"join tracks T ON PT.track_id=T.id join albums Al ON T.album_id=Al.id " +
			"join artists Ar ON Al.artist_id=Ar.id where PT.playlist_id = $2;",
			authID,
		playlistID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := &models.Track{}
		isFavourite := new(bool)
		isLiked := new(bool)

		if err := rows.Scan(&t.ID, &t.Name, &t.Duration, &t.Path, &t.Artist, &t.ArtistID, isFavourite, isLiked); err != nil {
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
