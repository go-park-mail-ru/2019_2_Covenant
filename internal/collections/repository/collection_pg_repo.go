package repository

import (
	"2019_2_Covenant/internal/collections"
	"2019_2_Covenant/internal/models"
	. "2019_2_Covenant/tools/vars"
	"database/sql"
)

type CollectionRepository struct {
	db *sql.DB
}

func NewCollectionRepository(db *sql.DB) collections.Repository {
	return &CollectionRepository{
		db: db,
	}
}

func (cr *CollectionRepository) Insert(collection *models.Collection) error {
	return cr.db.QueryRow("INSERT INTO collections (name, description) VALUES ($1, $2) RETURNING id, photo",
		collection.Name,
		collection.Description,
	).Scan(&collection.ID, &collection.Photo)
}

func (cr *CollectionRepository) DeleteByID(id uint64) error {
	if err := cr.db.QueryRow("DELETE FROM collections WHERE id = $1 RETURNING id",
		id,
	).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (cr *CollectionRepository) UpdateByID(collectionID uint64, name string, description string) error {
	if err := cr.db.QueryRow("UPDATE collections SET name = $1, description = $2 WHERE id = $3 RETURNING id",
		name,
		description,
		collectionID,
	).Scan(&collectionID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (cr *CollectionRepository) Select(count uint64, offset uint64) ([]*models.Collection, uint64, error) {
	var colls []*models.Collection
	var total uint64

	if err := cr.db.QueryRow("SELECT COUNT(*) FROM collections").Scan(&total); err != nil {
		return nil, total, err
	}

	rows, err := cr.db.Query("SELECT id, name, description, photo FROM collections " +
		"ORDER BY created_at LIMIT $1 OFFSET $2",
		count,
		offset,
	)

	if err != nil {
		return nil, total, err
	}

	defer rows.Close()

	for rows.Next() {
		c := &models.Collection{}

		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.Photo,
		); err != nil {
			return nil, total, err
		}

		colls = append(colls, c)
	}

	if err := rows.Err(); err != nil {
		return nil, total, err
	}

	return colls, total, nil
}

func (cr *CollectionRepository) SelectByID(id uint64) (*models.Collection, uint64, error) {
	c := &models.Collection{}
	var amountOfTracks uint64

	if err := cr.db.QueryRow("SELECT id, name, description, photo FROM collections WHERE id = $1",
		id,
	).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&c.Photo,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, amountOfTracks, ErrNotFound
		}

		return nil, amountOfTracks, err
	}

	if err := cr.db.QueryRow("SELECT COUNT(*) FROM collection_track WHERE collection_id = $1",
		id,
	).Scan(&amountOfTracks); err != nil {
		if err == sql.ErrNoRows {
			return c, amountOfTracks, nil
		}

		return nil, amountOfTracks, err
	}

	return c, amountOfTracks, nil
}

func (cr *CollectionRepository) InsertTrack(collectionID uint64, trackID uint64) error {
	if err := cr.db.QueryRow("SELECT id FROM tracks WHERE id = $1",
		trackID,
	).Scan(&trackID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	var id int
	if err := cr.db.QueryRow("SELECT id FROM collection_track WHERE collection_id = $1 AND track_id = $2",
		collectionID,
		trackID,
	).Scan(&id); err == nil {
		return ErrAlreadyExist
	}

	if _, err := cr.db.Exec("INSERT INTO collection_track (collection_id, track_id) VALUES ($1, $2)",
		collectionID,
		trackID,
	); err != nil {
		return err
	}

	return nil
}

func (cr *CollectionRepository) SelectTracks(collectionID uint64, authID uint64) ([]*models.Track, error) {
	var tracks []*models.Track

	rows, err := cr.db.Query(
		"select T.id, T.name, T.duration, T.path, Ar.name, Ar.id, " +
			"T.id in (select track_id from favourites where user_id = $1) AS favourite, " +
			"T.id in (select track_id from likes where user_id = $1) AS liked from collection_track CT " +
			"join tracks T ON CT.track_id=T.id join albums Al ON T.album_id=Al.id " +
			"join artists Ar ON Al.artist_id=Ar.id where CT.collection_id = $2;",
		authID,
		collectionID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := &models.Track{}
		isFavourite := new(bool)
		isLiked := new(bool)

		if err := rows.Scan(&t.ID, &t.Name, &t.Duration, &t.Path, &t.Artist,
			&t.ArtistID, isFavourite, isLiked); err != nil {
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

func (cr *CollectionRepository) UpdatePhoto(collectionID uint64, path string) error {
	if err := cr.db.QueryRow("UPDATE collections SET photo = $1 WHERE id = $2 RETURNING id",
		path,
		collectionID,
	).Scan(&collectionID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}

		return err
	}

	return nil
}
