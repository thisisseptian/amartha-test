package helper

import (
	"sync"

	"amartha-test/constant"
	"amartha-test/model"
)

var (
	userIDCounter int64
	mutexUser     sync.Mutex

	users = make(map[int64]*model.User)
)

func GenerateIncrementalUserID() int64 {
	mutexUser.Lock()
	defer mutexUser.Unlock()
	userIDCounter++
	return userIDCounter
}

func InitUsers() {
	var borrower1, lender1, lender2, fieldOfficer1 model.User

	borrower1 = model.User{
		UserID:   GenerateIncrementalUserID(),
		UserName: "Septian",
		UserType: constant.UserTypeBorrower,
	}

	lender1 = model.User{
		UserID:   GenerateIncrementalUserID(),
		UserName: "Pratama",
		UserType: constant.UserTypeLender,
	}

	lender2 = model.User{
		UserID:   GenerateIncrementalUserID(),
		UserName: "Rusmana",
		UserType: constant.UserTypeLender,
	}

	fieldOfficer1 = model.User{
		UserID:   GenerateIncrementalUserID(),
		UserName: "Amartha Officer",
		UserType: constant.UserTypeEmployee,
	}

	users[borrower1.UserID] = &borrower1
	users[lender1.UserID] = &lender1
	users[lender2.UserID] = &lender2
	users[fieldOfficer1.UserID] = &fieldOfficer1
}

func GetUsers() []model.User {
	var listUser []model.User
	for _, v := range users {
		listUser = append(listUser, *v)
	}

	return listUser
}

func GetUserByUserID(userID int64) model.User {
	user, exists := users[userID]
	if exists {
		return *user
	}

	return model.User{}
}
