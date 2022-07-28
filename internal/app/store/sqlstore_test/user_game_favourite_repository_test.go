package sqlstore_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store/sqlstore"
)

func TestUserGameFavouriteRepositoryFindByUserGame(t *testing.T) {
	userGameFavouriteWant := sqlstore.TestUserGameFavourites[1]

	userGameFavouriteFound, err := st.UserGameFavourites().FindByUserGame(userGameFavouriteWant.User, userGameFavouriteWant.Game)
	if err != nil {
		t.Errorf("Couldn't find userGameFavourite by userGame:\n\tUser: %+v\n\tGame: %+v\nError:\n\t%s", userGameFavouriteWant.User, userGameFavouriteWant.Game, err.Error())
	} else {
		if userGameFavouriteFound.ID != userGameFavouriteWant.ID || userGameFavouriteFound.Game.ID != userGameFavouriteWant.Game.ID || userGameFavouriteFound.User.ID != userGameFavouriteWant.User.ID {
			t.Errorf("Found wrong userGameFavourite by userGame:\n\tWanted: %+v,\n\tGot: %+v", userGameFavouriteWant, userGameFavouriteFound)
		}
	}
}

func TestUserGameFavouriteRepositoryDelete(t *testing.T) {
	userGameFavouriteWant := sqlstore.TestUserGameFavourites[2]

	userGameFavouriteFound, err := st.UserGameFavourites().Find(userGameFavouriteWant.ID)
	if err != nil {
		t.Errorf("Couldn't find userGameFavourite with ID (%d):\n\t%s", userGameFavouriteWant.ID, err.Error())
		return
	}
	if userGameFavouriteFound.ID != userGameFavouriteWant.ID || userGameFavouriteFound.Game.ID != userGameFavouriteWant.Game.ID || userGameFavouriteFound.User.ID != userGameFavouriteWant.User.ID {
		t.Errorf("Found wrong userGameFavourite by ID:\n\tWanted: %+v,\n\tGot: %+v", userGameFavouriteWant, userGameFavouriteFound)
	}

	if err := st.UserGameFavourites().Delete(userGameFavouriteWant.ID); err != nil {
		t.Errorf("Couldn't delete userGameFavourite with ID (%d):\n\t%s", userGameFavouriteWant.ID, err.Error())
	}

	if userGameFavouriteNotExist, err := st.UserGameFavourites().Find(userGameFavouriteWant.ID); err == nil {
		t.Errorf("Found user with non-existent ID (%d):\n\t%+v", userGameFavouriteWant.ID, userGameFavouriteNotExist)
	} else {
		if errors.Cause(err) != store.ErrNotFound {
			t.Errorf("Wrong error when finding user with non-existent ID (%d):\n\t%s", userGameFavouriteWant.ID, err.Error())
		}
	}
}
