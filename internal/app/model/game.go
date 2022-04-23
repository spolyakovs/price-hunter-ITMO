package model

type Game struct {
	ID             uint64     `json:"id" db:"id,omitempty"`
	HeaderImageURL string     `json:"header_image" db:"header_image_url"`
	Name           string     `json:"name" db:"name"`
	Description    string     `json:"description" db:"description"`
	ReleaseDate    string     `json:"release_date" db:"release_date"` // format "dd.MM.YYYY"
	Publisher      *Publisher `json:"publisher" db:"publisher"`
}
