package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var userContextKey contextKey = contextKey("user")

func (h *Handlers) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromCookie(r, h.store.User())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
	})
}

func GetUser(r *http.Request) syncedvideo.User {
	return r.Context().Value(userContextKey).(syncedvideo.User)
}

const userCookieKey string = "userID"

func hasUserCookie(r *http.Request) bool {
	c, err := r.Cookie(userCookieKey)
	if err != nil {
		return false
	}
	return c.Value != ""
}

func getUserFromCookie(r *http.Request, userStore syncedvideo.UserStore) (syncedvideo.User, error) {
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
	user, err := userStore.Get(userID)
	if err != nil {
		return syncedvideo.User{}, err
	}
	return user, nil
}

func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	user := syncedvideo.User{}
	if hasUserCookie(r) {
		u, err := getUserFromCookie(r, h.store.User())
		if err != nil && err != sql.ErrNoRows {
			log.Printf("error getting user: %v", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		user = u
	}

	if user.ID == uuid.Nil {
		err := h.store.User().Create(&user)
		if err != nil {
			log.Printf("error creating user: %v", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     userCookieKey,
		Value:    user.ID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().UTC().Add(24 * time.Hour * 30), // 30 days
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
