package model

type Game struct {
	ID             uint64     `json:"id" db:"id,omitempty"`
	HeaderImageURL string     `json:"header_image" db:"header_image_url"`
	Name           string     `json:"name" db:"name"`
	Description    string     `json:"description" db:"description"`
	Publisher      *Publisher `json:"publisher" db:"publisher"`
}
