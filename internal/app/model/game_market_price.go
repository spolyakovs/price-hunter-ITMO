package model

type GameMarketPrice struct {
	ID                    uint64  `json:"id" db:"id,omitempty"`
	InitialValueFormatted string  `json:"initial_value_formatted" db:"initial_value_formatted"`
	FinalValueFormatted   string  `json:"final_value_formatted" db:"final_value_formatted"`
	DiscountPercent       int     `json:"discount_percent" db:"discount_percent"`
	MarketGameURL         string  `json:"uri_string" db:"market_game_url"`
	Game                  *Game   `json:"game" db:"game"`
	Market                *Market `json:"market" db:"market"`
}
