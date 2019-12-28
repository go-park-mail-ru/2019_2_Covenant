package models

type Collection struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
}

func NewCollection(name string, description string) *Collection {
	return &Collection{
		Name:        name,
		Description: description,
	}
}
