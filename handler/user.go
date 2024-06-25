package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"amartha-test/helper"
)

// ListUser ...
func ListUser(w http.ResponseWriter, r *http.Request) {
	// 1. check http method
	if r.Method != http.MethodGet {
		log.Println("[ListUser] invalid request method")
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 2. get user list
	users := helper.GetUsers()
	if len(users) == 0 {
		log.Println("[ListUser] user list is empty")
		http.Error(w, "user list is empty", http.StatusNotFound)
		return
	}

	// 3. return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
