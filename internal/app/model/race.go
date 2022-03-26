package model

type Race struct {
	ID       int    `json:"id" db:"id,omitempty"`
	Name     string `json:"name" db:"name"`
	Location string `json:"location" db:"location"`
	Date     string `json:"date" db:"date"` // format "YYYY-MM-dd"
}
