package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/syncedvideo/syncedvideo"
)

type UserHandler struct {
	store syncedvideo.Store
}

func NewUserHandler(s syncedvideo.Store) syncedvideo.UserHandler {
	return &UserHandler{
		store: s,
	}
}

func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	// get userID from userID cookie
	userIDCookie, err := r.Cookie("userID")
	userID := uuid.New()
	if err == nil {
		userID, err = uuid.Parse(userIDCookie.Value)
		if err != nil {
			log.Printf("error parsing uuid: %v", err)
			return
		}
	}

	// get or create user
	user, err := h.store.User().Get(userID)
	if err == sql.ErrNoRows {
		user.ID = userID
		err2 := h.store.User().Create(&user)
		if err2 != nil {
			log.Println(err2)
			return
		}
	} else if err != nil {
		log.Println(err)
		return
	}

	// set userID cookie
	userIDCookie = &http.Cookie{
		Name:     "userID",
		Value:    user.ID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().UTC().Add(24 * time.Hour * 30), // 30 days
	}
	http.SetCookie(w, userIDCookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
