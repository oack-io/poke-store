package middleware

import (
	"context"
	"net/http"

	"github.com/oack-io/poke-store/internal/model"
	"github.com/oack-io/poke-store/internal/store"
)

type contextKey string

const userKey contextKey = "user"

func Auth(sessions *store.SessionStore) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"not authenticated"}`))
				return
			}

			user, ok := sessions.Get(cookie.Value)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"session expired"}`))
				return
			}

			ctx := context.WithValue(r.Context(), userKey, user)
			next(w, r.WithContext(ctx))
		}
	}
}

func UserFromContext(ctx context.Context) model.User {
	user, _ := ctx.Value(userKey).(model.User)
	return user
}
