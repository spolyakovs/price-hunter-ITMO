package apiStore

import "github.com/pkg/errors"

var (
	ErrAppSkiped = errors.New("App skipped")
	// ErrGameInfo = errors.New("Couldn't get game info")
)

const (
	errAPIStoreMessageFormat = "API %s method %s error"
)
