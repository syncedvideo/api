package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/syncedvideo/syncedvideo"
	"github.com/syncedvideo/syncedvideo/http/middleware"
	"github.com/syncedvideo/syncedvideo/http/request"
	"github.com/syncedvideo/syncedvideo/http/response"
)

type userHandler struct{}

func RegisterUserHandler(r chi.Router) {
	userHandler := newUserHandler()
	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.UserMiddleware)
		r.Put("/", userHandler.Update)
	})
}

func newUserHandler() *userHandler {
	return &userHandler{}
}

const userCookieKey string = "userID"

func hasUserCookie(r *http.Request) bool {
	c, err := r.Cookie(userCookieKey)
	if err != nil {
		return false
	}
	return c.Value != ""
}

func getUserFromCookie(r *http.Request) (syncedvideo.User, error) {
	userIDCookie, err := r.Cookie(userCookieKey)
	if err != nil {
		return syncedvideo.User{}, err
	}
	if userIDCookie.Value == "" {
		return syncedvideo.User{}, errors.New(userCookieKey + " cookie value is empty")
	}

	userID, err := uuid.Parse(userIDCookie.Value)
	if err != nil {
		return syncedvideo.User{}, err
	}
	if userID == uuid.Nil {
		return syncedvideo.User{}, errors.New("userID is nil")
	}

	user, err := syncedvideo.Config.Store.User().Get(userID)
	if err != nil {
		return syncedvideo.User{}, err
	}
	return user, nil
}

func (h *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	currentUser := request.GetUserCtx(r)
	updatedUser := syncedvideo.User{}
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		log.Printf("Decode failed: %v", err)
		response.WithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	if currentUser.ID != updatedUser.ID {
		log.Printf("Update failed: user id %s tried to update user id %s", currentUser.ID, updatedUser.ID)
		response.WithError(w, "cannot update different user", http.StatusForbidden)
		return
	}
	err = syncedvideo.Config.Store.User().Update(&currentUser)
	if err != nil {
		log.Printf("Update failed: %v", err)
		response.WithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	rooms, err := syncedvideo.Config.Store.User().GetCurrentRooms(currentUser.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("GetCurrentRooms failed: %v", err)
		response.WithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	for _, room := range rooms {
		err := room.SyncUsers()
		if err != nil {
			log.Printf("SyncUsers failed: %v", err)
			response.WithError(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
	response.WithJSON(w, updatedUser, http.StatusOK)
}
