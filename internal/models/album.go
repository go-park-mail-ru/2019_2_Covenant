package models

type Album struct {
	ID       uint64    `json:"id"`
	ArtistID uint64    `json:"artist_id"`
	Name     string    `json:"name"`
	Photo    string    `json:"photo"`
	Year     string    `json:"year"`
}

func NewAlbum(name string, year string, artistID uint64) *Album {
	return &Album{
		ArtistID: artistID,
		Name: name,
		Year: year,
	}
}
