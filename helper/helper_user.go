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

func (h *Helper) GenerateIncrementalUserID() int64 {
	mutexUser.Lock()
	defer mutexUser.Unlock()
	userIDCounter++
	return userIDCounter
}

func (h *Helper) InitUsers() {
	var borrower1, lender1, lender2, fieldValidator1, fieldOfficer1 model.User

	borrower1 = model.User{
		UserID:   h.GenerateIncrementalUserID(),
		UserName: "Septian",
		UserType: constant.UserTypeBorrower,
	}

	lender1 = model.User{
		UserID:   h.GenerateIncrementalUserID(),
		UserName: "Pratama",
		UserType: constant.UserTypeLender,
	}

	lender2 = model.User{
		UserID:   h.GenerateIncrementalUserID(),
		UserName: "Rusmana",
		UserType: constant.UserTypeLender,
	}

	fieldValidator1 = model.User{
		UserID:   h.GenerateIncrementalUserID(),
		UserName: "Validator",
		UserType: constant.UserTypeFieldValidatorEmployee,
	}

	fieldOfficer1 = model.User{
		UserID:   h.GenerateIncrementalUserID(),
		UserName: "Officer",
		UserType: constant.UserTypeFieldOfficerEmployee,
	}

	users[borrower1.UserID] = &borrower1
	users[lender1.UserID] = &lender1
	users[lender2.UserID] = &lender2
	users[fieldValidator1.UserID] = &fieldValidator1
	users[fieldOfficer1.UserID] = &fieldOfficer1
}

func (h *Helper) GetUsers() []model.User {
	var listUser []model.User
	for _, v := range users {
		listUser = append(listUser, *v)
	}

	return listUser
}

func (h *Helper) GetUserByUserID(userID int64) model.User {
	user, exists := users[userID]
	if exists {
		return *user
	}

	return model.User{}
}
