package handler

import (
	"encoding/json"
	"net/http"

	"github.com/oack-io/poke-store/internal/data"
	"github.com/oack-io/poke-store/internal/middleware"
	"github.com/oack-io/poke-store/internal/model"
	"github.com/oack-io/poke-store/internal/store"
)

type Cart struct {
	carts   *store.CartStore
	catalog *data.Catalog
}

func NewCart(carts *store.CartStore, catalog *data.Catalog) *Cart {
	return &Cart{carts: carts, catalog: catalog}
}

func (h *Cart) Get(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	cart := h.carts.Get(user.Email)

	detail := h.buildCartDetail(cart)
	JSON(w, http.StatusOK, detail)
}

func (h *Cart) Add(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())

	var req struct {
		PokemonID int `json:"pokemonId"`
		Quantity  int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	if _, ok := h.catalog.Get(req.PokemonID); !ok {
		Error(w, http.StatusNotFound, "pokemon not found")
		return
	}

	h.carts.Add(user.Email, req.PokemonID, req.Quantity)

	cart := h.carts.Get(user.Email)
	detail := h.buildCartDetail(cart)
	JSON(w, http.StatusOK, detail)
}

func (h *Cart) Remove(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())

	var req struct {
		PokemonID int `json:"pokemonId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	h.carts.Remove(user.Email, req.PokemonID)

	cart := h.carts.Get(user.Email)
	detail := h.buildCartDetail(cart)
	JSON(w, http.StatusOK, detail)
}

func (h *Cart) Clear(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	h.carts.Clear(user.Email)
	JSON(w, http.StatusOK, map[string]any{
		"items": []any{},
		"total": 0,
	})
}

func (h *Cart) buildCartDetail(cart *model.Cart) CartDetailResponse {
	var items []cartItemResponse
	total := 0

	for _, item := range cart.Items {
		pokemon, ok := h.catalog.Get(item.PokemonID)
		if !ok {
			continue
		}
		subtotal := pokemon.Price * item.Quantity
		total += subtotal
		items = append(items, cartItemResponse{
			Pokemon:  pokemon,
			Quantity: item.Quantity,
			Subtotal: subtotal,
		})
	}

	if items == nil {
		items = []cartItemResponse{}
	}

	return CartDetailResponse{Items: items, Total: total}
}

type cartItemResponse struct {
	Pokemon  any `json:"pokemon"`
	Quantity int `json:"quantity"`
	Subtotal int `json:"subtotal"`
}

type CartDetailResponse struct {
	Items []cartItemResponse `json:"items"`
	Total int                `json:"total"`
}
