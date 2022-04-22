package model

type GameTag struct {
	ID   uint64 `json:"id" db:"id,omitempty"`
	Game *Game  `json:"game" db:"game"`
	Tag  *Tag   `json:"tag" db:"tag"`
}
