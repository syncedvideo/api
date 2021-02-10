package handler

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/syncedvideo/syncedvideo"
	"github.com/syncedvideo/syncedvideo/http/response"
)

type authHandler struct{}

func RegisterAuthHandler(r chi.Router) {
	authHandler := newAuthHandler()
	r.Post("/auth", authHandler.Auth)
}

func newAuthHandler() *authHandler {
	return &authHandler{}
}

func (h *authHandler) Auth(w http.ResponseWriter, r *http.Request) {
	user := syncedvideo.User{}

	if hasUserCookie(r) {
		u, err := getUserFromCookie(r)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("error getting user: %v", err)
			response.WithError(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		user = u
	}

	if user.ID == uuid.Nil {
		user = syncedvideo.NewUser()
		err := syncedvideo.Config.Store.User().Create(&user)
		if err != nil {
			log.Printf("error creating user: %v", err)
			response.WithError(w, "something went wrong", http.StatusInternalServerError)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     userCookieKey,
		Value:    user.ID.String(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Expires:  time.Now().UTC().Add(24 * time.Hour * 30), // 30 days
	})

	response.WithJSON(w, user, http.StatusOK)
}
