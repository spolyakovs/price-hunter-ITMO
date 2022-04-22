package model

type Tag struct {
	ID   uint64 `json:"id" db:"id,omitempty"`
	Name string `json:"name" db:"name"`
}
