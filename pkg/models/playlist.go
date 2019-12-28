package models

type Playlist struct {
	ID            uint64 `json:"id"`
	OwnerID       uint64 `json:"owner_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Photo         string `json:"photo"`
}

func NewPlaylist(name string, description string, ownerID uint64) *Playlist {
	return &Playlist{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}
}
