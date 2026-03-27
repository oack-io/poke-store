package model

type Pokemon struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Type        []string `json:"type"`
	Description string   `json:"description"`
	Price       int      `json:"price"` // in PokeDollars (cents)
	Image       string   `json:"image"`
	HP          int      `json:"hp"`
	Attack      int      `json:"attack"`
	Defense     int      `json:"defense"`
	Speed       int      `json:"speed"`
}

type CartItem struct {
	PokemonID int `json:"pokemonId"`
	Quantity  int `json:"quantity"`
}

type Cart struct {
	Items []CartItem `json:"items"`
}

type CartItemDetail struct {
	Pokemon  Pokemon `json:"pokemon"`
	Quantity int     `json:"quantity"`
	Subtotal int     `json:"subtotal"`
}

type CartDetail struct {
	Items []CartItemDetail `json:"items"`
	Total int              `json:"total"`
}

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

type Order struct {
	ID     string           `json:"id"`
	User   string           `json:"user"`
	Items  []CartItemDetail `json:"items"`
	Total  int              `json:"total"`
	Status string           `json:"status"`
}
