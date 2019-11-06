package models

type Track struct {
	ID       uint64        `json:"-"`
	AlbumID  uint64        `json:"-"`
	ArtistID uint64		   `json:"-"`
	Name     string        `json:"name"`
	Duration string        `json:"duration"`
	Photo    string        `json:"photo"`
	Artist   string        `json:"artist"`
	Album    string        `json:"album"`
	Path     string        `json:"path"`
}
