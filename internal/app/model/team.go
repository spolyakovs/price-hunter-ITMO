package model

type Team struct {
	ID                 int    `json:"id" db:"id,omitempty"`
	Name               string `json:"name" db:"name"`
	EngineManufacturer string `json:"engine_manufacturer" db:"engine_manufacturer"`
}
