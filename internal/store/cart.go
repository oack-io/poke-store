package store

import (
	"sync"

	"github.com/oack-io/poke-store/internal/model"
)

type CartStore struct {
	mu    sync.RWMutex
	carts map[string]*model.Cart // keyed by user email
}

func NewCartStore() *CartStore {
	return &CartStore{
		carts: make(map[string]*model.Cart),
	}
}

func (s *CartStore) Get(email string) *model.Cart {
	s.mu.RLock()
	cart, ok := s.carts[email]
	s.mu.RUnlock()
	if !ok {
		return &model.Cart{Items: []model.CartItem{}}
	}
	return cart
}

func (s *CartStore) Add(email string, pokemonID int, quantity int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.carts[email]
	if !ok {
		cart = &model.Cart{Items: []model.CartItem{}}
		s.carts[email] = cart
	}

	for i, item := range cart.Items {
		if item.PokemonID == pokemonID {
			cart.Items[i].Quantity += quantity
			return
		}
	}
	cart.Items = append(cart.Items, model.CartItem{PokemonID: pokemonID, Quantity: quantity})
}

func (s *CartStore) Remove(email string, pokemonID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.carts[email]
	if !ok {
		return
	}

	for i, item := range cart.Items {
		if item.PokemonID == pokemonID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			return
		}
	}
}

func (s *CartStore) Clear(email string) {
	s.mu.Lock()
	delete(s.carts, email)
	s.mu.Unlock()
}
