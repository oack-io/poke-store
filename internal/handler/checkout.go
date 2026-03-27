package handler

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"

	"github.com/oack-io/poke-store/internal/data"
	"github.com/oack-io/poke-store/internal/middleware"
	"github.com/oack-io/poke-store/internal/model"
	"github.com/oack-io/poke-store/internal/store"
)

type Checkout struct {
	carts   *store.CartStore
	catalog *data.Catalog
}

func NewCheckout(carts *store.CartStore, catalog *data.Catalog) *Checkout {
	return &Checkout{carts: carts, catalog: catalog}
}

func (h *Checkout) Process(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	cart := h.carts.Get(user.Email)

	if len(cart.Items) == 0 {
		Error(w, http.StatusBadRequest, "cart is empty")
		return
	}

	var items []model.CartItemDetail
	total := 0
	for _, item := range cart.Items {
		pokemon, ok := h.catalog.Get(item.PokemonID)
		if !ok {
			continue
		}
		subtotal := pokemon.Price * item.Quantity
		total += subtotal
		items = append(items, model.CartItemDetail{
			Pokemon:  pokemon,
			Quantity: item.Quantity,
			Subtotal: subtotal,
		})
	}

	orderID := generateOrderID()
	order := model.Order{
		ID:     orderID,
		User:   user.Email,
		Items:  items,
		Total:  total,
		Status: "confirmed",
	}

	h.carts.Clear(user.Email)

	slog.Info("order placed",
		"scope", "cart",
		"orderID", orderID,
		"user", user.Email,
		"total", total,
		"itemCount", len(items),
	)

	JSON(w, http.StatusOK, order)
}

func generateOrderID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "ORD-" + hex.EncodeToString(b)
}
