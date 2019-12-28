package models

type Track struct {
	ID          uint64 `json:"id"`
	AlbumID     uint64 `json:"album_id,omitempty"`
	ArtistID    uint64 `json:"artist_id,omitempty"`
	Name        string `json:"name"`
	Duration    string `json:"duration"`
	Photo       string `json:"photo,omitempty"`
	Artist      string `json:"artist,omitempty"`
	Album       string `json:"album"`
	Path        string `json:"path"`
	IsFavourite *bool  `json:"is_favourite,omitempty"`
	IsLiked     *bool  `json:"is_liked,omitempty"`
}
