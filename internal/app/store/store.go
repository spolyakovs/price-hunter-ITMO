package store

// TODO: change tables structure for Price Hunter app

type Store interface {
	Users() UserRepository
	Publishers() PublisherRepository
}
