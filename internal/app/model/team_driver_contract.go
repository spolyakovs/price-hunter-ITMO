package model

type TeamDriverContract struct {
	ID       int     `json:"id" db:"id,omitempty"`
	FromDate string  `json:"from_date" db:"from_date"`                 // format "YYYY-MM-dd"
	ToDate   string  `json:"to_date,omitempty" db:"to_date,omitempty"` // format "YYYY-MM-dd"
	Team     *Team   `json:"team" db:"team"`
	Driver   *Driver `json:"driver" db:"driver"`
}
