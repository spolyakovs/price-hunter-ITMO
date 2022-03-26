package model

type RaceResult struct {
	ID     int     `json:"id" db:"id,omitempty"`
	Place  int     `json:"place" db:"place"`
	Points int     `json:"points" db:"points"`
	Race   *Race   `json:"race" db:"race"`
	Driver *Driver `json:"driver" db:"driver"`
}
