package model

type Driver struct {
	ID        int    `json:"id" db:"id,omitempty"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`   // format "YYYY-MM-dd"
	BirthDate string `json:"birth_date" db:"birth_date"` // format "YYYY-MM-dd"
}
