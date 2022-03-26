package store

// TODO: change tables structure for Price Hunter app

type Store interface {
	Users() UserRepository
	Teams() TeamRepository
	Drivers() DriverRepository
	Races() RaceRepository
	TeamDriverContracts() TeamDriverContractRepository
	RaceResults() RaceResultRepository
}
