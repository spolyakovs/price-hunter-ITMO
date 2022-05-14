package sqlstore_test

import (
	"testing"
)

func TestUserRepositoryCreate(t *testing.T) {
	store, err := setupStore()
	if  err != nil {
		t.Error(err.Error())
		return
	}

	if err := store.ClearTables(); err != nil {
		t.Errorf("Couldn't clear test data from DB:\n\t%s", err.Error())
	}
}
