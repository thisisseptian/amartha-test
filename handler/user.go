package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ListUser is handler to get list of users
func (h *Handler) ListUser(w http.ResponseWriter, r *http.Request) {
	// 1. get user list
	users := h.Helper.GetUsers()
	if len(users) == 0 {
		log.Println("[ListUser] user list is empty")
		h.RenderResponse(w, r, "", http.StatusNotFound, "[ListUser] user list is empty")
		return
	}

	// 2. render response
	h.RenderResponse(w, r, users, http.StatusOK, "")
}

// DetailUser is handler to get user detail
func (h *Handler) DetailUser(w http.ResponseWriter, r *http.Request) {
	// 1. get vars
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		log.Printf("[DetailUser] failed parse int, with error: %+v", err)
		h.RenderResponse(w, r, "", http.StatusBadRequest, fmt.Sprintf("[DetailUser] failed parse int, with error: %+v", err))
		return
	}

	// 2. sanitize payload
	if userID == 0 {
		log.Println("[DetailUser] user id is zero")
		h.RenderResponse(w, r, "", http.StatusBadRequest, "[DetailUser] user id is zero")
		return
	}

	// 3. get user by user id
	user := h.Helper.GetUserByUserID(userID)
	if user.UserID == 0 {
		log.Printf("[DetailUser][UserID: %d] user data is not found", userID)
		h.RenderResponse(w, r, "", http.StatusNotFound, fmt.Sprintf("[DetailUser][UserID: %d] user data is not found", userID))
		return
	}

	// 4. render response
	h.RenderResponse(w, r, user, http.StatusOK, "")
}
