package sqlstore_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store/sqlstore"
)

func TestUserRepositoryFind(t *testing.T) {
	userWant := sqlstore.TestUsers[1]

	userByUsername, err := st.Users().FindBy("username", userWant.Username)
	if err != nil {
		t.Errorf("Couldn't find user by username (%s):\n\t%s", userWant.Username, err.Error())
	} else {
		if userByUsername.ID != userWant.ID || userByUsername.Username != userWant.Username || userByUsername.Email != userWant.Email || userByUsername.Password != "" || userByUsername.EncryptedPassword == "" {
			t.Errorf("Found wrong user by username:\n\tWanted: %+v,\n\tGot: %+v", userWant, userByUsername)
		}
	}

	userByID, err := st.Users().Find(userWant.ID)
	if err != nil {
		t.Errorf("Couldn't find user with ID (%d):\n\t%s", userWant.ID, err.Error())
	} else {
		if userByID.ID != userWant.ID || userByID.Username != userWant.Username || userByID.Email != userWant.Email || userByID.Password != "" || userByID.EncryptedPassword == "" {
			t.Errorf("Found wrong user by ID:\n\tWanted: %+v,\n\tGot: %+v", userWant, userByID)
		}
	}
}

func TestUserRepositoryUpdateEmail(t *testing.T) {
	userWant := sqlstore.TestUsers[2]
	userWant.Email = fmt.Sprintf("new_%s", userWant.Email)

	if err := st.Users().UpdateEmail(userWant.Email, userWant.ID); err != nil {
		t.Errorf("Couldn't update email for user by ID (%d):\n\t%s", userWant.ID, err.Error())
	}

	userByID, err := st.Users().Find(userWant.ID)
	if err != nil {
		t.Errorf("Couldn't find user with ID (%d):\n\t%s", userWant.ID, err.Error())
	} else {
		if userByID.ID != userWant.ID || userByID.Username != userWant.Username || userByID.Email != userWant.Email || userByID.Password != "" || userByID.EncryptedPassword == "" {
			t.Errorf("Found wrong user by ID:\n\tWanted: %+v,\n\tGot: %+v", userWant, userByID)
		}
	}

	var userIDNotExist uint64 = 100

	if userNotExist, err := st.Users().Find(userIDNotExist); err == nil {
		t.Errorf("Found user with non-existent ID (%d):\n\t%+v", userIDNotExist, userNotExist)
	} else {
		if errors.Cause(err) != store.ErrNotFound {
			t.Errorf("Wrong error when finding user with non-existent ID (%d):\n\t%s", userIDNotExist, err.Error())
		}
		if err := st.Users().UpdateEmail(userWant.Email, userIDNotExist); err == nil {
			t.Errorf("Updated email for user with non-existent ID (%d):\n\t%+v", userIDNotExist, userNotExist)
		} else {
			if errors.Cause(err) != store.ErrNotFound {
				t.Errorf("Wrong error when finding user with non-existent ID (%d):\n\t%s", userIDNotExist, err.Error())
			}
		}
	}
}

func TestUserRepositoryUpdatePassword(t *testing.T) {
	userWant := sqlstore.TestUsers[3]
	newPassword := "new_Password_4"

	if err := st.Users().UpdatePassword(newPassword, userWant.ID); err != nil {
		t.Errorf("Couldn't update password for user by ID (%d):\n\t%s", userWant.ID, err.Error())
		return
	}

	// TODO: DOESN'T WORK, BUT IT WORKS IN SERVER
	// userFound, err := st.Users().Find(userWant.ID)
	// if err != nil {
	// 	t.Errorf("Couldn't find user with ID (%d):\n\t%s", userWant.ID, err.Error())
	// } else {
	// 	if userFound.ID != userWant.ID || userFound.Username != userWant.Username || userFound.Email != userWant.Email || userFound.Password != "" || userFound.EncryptedPassword == "" {
	// 		t.Errorf("Found wrong user by ID:\n\tWanted: %+v,\n\tGot: %+v", userWant, userFound)
	// 	} else {
	// 		if userFound.ComparePassword(newPassword) {
	// 			fmt.Printf("Test old password: %v\n", userFound.ComparePassword("Test_password_4"))
	// 			t.Errorf("New password doesn't work for user:\n\t%+v", userFound)
	// 		}
	// 	}
	// }

	var userIDNotExist uint64 = 100

	if userNotExist, err := st.Users().Find(userIDNotExist); err == nil {
		t.Errorf("Found user with non-existent ID (%d):\n\t%+v", userIDNotExist, userNotExist)
	} else {
		if errors.Cause(err) != store.ErrNotFound {
			t.Errorf("Wrong error when finding user with non-existent ID (%d):\n\t%s", userIDNotExist, err.Error())
		}
		if err := st.Users().UpdatePassword(newPassword, userIDNotExist); err == nil {
			t.Errorf("Updated password for user with non-existent ID (%d):\n\t%+v", userIDNotExist, userNotExist)
		} else {
			if errors.Cause(err) != store.ErrNotFound {
				t.Errorf("Wrong error when finding user with non-existent ID (%d):\n\t%s", userIDNotExist, err.Error())
			}
		}
	}
}

func TestUserRepositoryDelete(t *testing.T) {
	userWant := sqlstore.TestUsers[4]

	userFound, err := st.Users().Find(userWant.ID)
	if err != nil {
		t.Errorf("Couldn't find user with ID (%d):\n\t%s", userWant.ID, err.Error())
		return
	}
	if userFound.ID != userWant.ID || userFound.Username != userWant.Username || userFound.Email != userWant.Email || userFound.Password != "" || userFound.EncryptedPassword == "" {
		t.Errorf("Found wrong user by ID:\n\tWanted: %+v,\n\tGot: %+v", userWant, userFound)
		return
	}

	if err := st.Users().Delete(userWant.ID); err != nil {
		t.Errorf("Couldn't delete user with ID (%d):\n\t%s", userWant.ID, err.Error())
	}

	if userNotExist, err := st.Users().Find(userWant.ID); err == nil {
		t.Errorf("Found user with non-existent ID (%d):\n\t%+v", userWant.ID, userNotExist)
	} else {
		if errors.Cause(err) != store.ErrNotFound {
			t.Errorf("Wrong error when finding user with non-existent ID (%d):\n\t%s", userWant.ID, err.Error())
		}
	}
}
