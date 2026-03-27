package data

import "github.com/oack-io/poke-store/internal/model"

// DemoUsers returns the hardcoded demo user accounts.
// Passwords are stored in plain text because this is a demo app for Playwright testing.
var DemoUsers = []model.User{
	{Email: "ash@pokemon.com", Name: "Ash Ketchum", Password: "pikachu123"},
	{Email: "misty@pokemon.com", Name: "Misty Waterflower", Password: "starmie123"},
	{Email: "brock@pokemon.com", Name: "Brock Harrison", Password: "onix123"},
}

func FindUser(email, password string) (model.User, bool) {
	for _, u := range DemoUsers {
		if u.Email == email && u.Password == password {
			return u, true
		}
	}
	return model.User{}, false
}
