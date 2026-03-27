package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/oack-io/poke-store/internal/data"
	"github.com/oack-io/poke-store/internal/middleware"
	"github.com/oack-io/poke-store/internal/store"
)

type Auth struct {
	sessions *store.SessionStore
	carts    *store.CartStore
}

func NewAuth(sessions *store.SessionStore, carts *store.CartStore) *Auth {
	return &Auth{sessions: sessions, carts: carts}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, ok := data.FindUser(req.Email, req.Password)
	if !ok {
		slog.Info("login failed", "scope", "auth", "email", req.Email)
		Error(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token := h.sessions.Create(user)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	slog.Info("login success", "scope", "auth", "email", user.Email)
	JSON(w, http.StatusOK, map[string]any{
		"user": map[string]string{
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		h.sessions.Delete(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Auth) Me(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	cart := h.carts.Get(user.Email)
	cartCount := 0
	for _, item := range cart.Items {
		cartCount += item.Quantity
	}

	JSON(w, http.StatusOK, map[string]any{
		"email":     user.Email,
		"name":      user.Name,
		"cartCount": cartCount,
	})
}
