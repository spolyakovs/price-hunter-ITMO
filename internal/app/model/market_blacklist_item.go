package model

type MarketBlacklistItem struct {
	ID            uint64  `db:"id,omitempty"`
	MarketGameURL string  `db:"market_game_url"`
	Market        *Market `db:"market"`
}
