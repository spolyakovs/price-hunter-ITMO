package model

type UserGameFavourite struct {
	ID   uint64 `json:"id" db:"id,omitempty"`
	User *User  `json:"user" db:"user"`
	Game *Game  `json:"game" db:"game"`
}
