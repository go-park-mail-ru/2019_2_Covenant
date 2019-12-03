package models

type Artist struct {
	ID    uint64 `json:"-"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}
