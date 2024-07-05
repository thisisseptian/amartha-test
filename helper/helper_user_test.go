package helper

import (
	"testing"
)

func TestInitUsersAndGetUsers(t *testing.T) {
	helper := NewHelper()

	t.Run("initialize and get users", func(t *testing.T) {
		helper.InitUsers()

		usersList := helper.GetUsers()
		if len(usersList) != 5 {
			t.Errorf("expected 5 users, got %d", len(usersList))
		}

		expectedUsernames := []string{"Septian", "Pratama", "Rusmana", "Validator", "Officer"}
		for _, username := range expectedUsernames {
			found := false
			for _, user := range usersList {
				if user.UserName == username {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected user %s not found", username)
			}
		}
	})
}

func TestGetUserByUserID(t *testing.T) {
	helper := NewHelper()

	t.Run("get user by user id", func(t *testing.T) {
		helper.InitUsers()

		userID := int64(1)
		user := helper.GetUserByUserID(userID)
		if user.UserName != "Septian" && user.UserName != "Pratama" && user.UserName != "Rusmana" && user.UserName != "Validator" && user.UserName != "Officer" {
			userID = int64(7)
			user = helper.GetUserByUserID(userID)
			if user.UserName != "Septian" && user.UserName != "Pratama" && user.UserName != "Rusmana" && user.UserName != "Validator" && user.UserName != "Officer" {
				t.Errorf("expected username 'Septian', got '%s'", user.UserName)
			}
		}

		userIDNotFound := int64(999)
		userNotFound := helper.GetUserByUserID(userIDNotFound)
		if userNotFound.UserID != 0 {
			t.Errorf("expected non-existing user to return empty user, got %+v", userNotFound)
		}
	})
}
