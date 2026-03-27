package data

import (
	"strings"

	"github.com/oack-io/poke-store/internal/model"
)

type Catalog struct {
	pokemon []model.Pokemon
	byID    map[int]model.Pokemon
}

func NewCatalog() *Catalog {
	c := &Catalog{
		pokemon: pokemonList(),
		byID:    make(map[int]model.Pokemon),
	}
	for _, p := range c.pokemon {
		c.byID[p.ID] = p
	}
	return c
}

func (c *Catalog) All() []model.Pokemon {
	return c.pokemon
}

func (c *Catalog) Search(query string) []model.Pokemon {
	if query == "" {
		return c.pokemon
	}
	q := strings.ToLower(query)
	var results []model.Pokemon
	for _, p := range c.pokemon {
		if strings.Contains(strings.ToLower(p.Name), q) {
			results = append(results, p)
			continue
		}
		for _, t := range p.Type {
			if strings.Contains(strings.ToLower(t), q) {
				results = append(results, p)
				break
			}
		}
	}
	return results
}

func (c *Catalog) Get(id int) (model.Pokemon, bool) {
	p, ok := c.byID[id]
	return p, ok
}

func pokemonList() []model.Pokemon {
	return []model.Pokemon{
		{
			ID: 1, Name: "Bulbasaur", Type: []string{"Grass", "Poison"},
			Description: "A strange seed was planted on its back at birth. The plant sprouts and grows with this Pokemon. Perfect as a starter companion or a lovely desk plant that occasionally uses Vine Whip on your coworkers.",
			Price: 2999, Image: "/pokemon/bulbasaur.png", HP: 45, Attack: 49, Defense: 49, Speed: 45,
		},
		{
			ID: 4, Name: "Charmander", Type: []string{"Fire"},
			Description: "The flame on its tail indicates its life force. If it is healthy, the flame burns brightly. Great for lighting birthday candles, campfires, and your rival's ego.",
			Price: 3499, Image: "/pokemon/charmander.png", HP: 39, Attack: 52, Defense: 43, Speed: 65,
		},
		{
			ID: 7, Name: "Squirtle", Type: []string{"Water"},
			Description: "When it retracts its long neck into its shell, it squirts out water with vigorous force. Also doubles as the world's cutest water gun at pool parties.",
			Price: 2999, Image: "/pokemon/squirtle.png", HP: 44, Attack: 48, Defense: 65, Speed: 43,
		},
		{
			ID: 25, Name: "Pikachu", Type: []string{"Electric"},
			Description: "When several of these Pokemon gather, their electricity can build and cause lightning storms. WARNING: Do not use as a phone charger. We learned that the hard way.",
			Price: 8999, Image: "/pokemon/pikachu.png", HP: 35, Attack: 55, Defense: 40, Speed: 90,
		},
		{
			ID: 6, Name: "Charizard", Type: []string{"Fire", "Flying"},
			Description: "It spits fire that is hot enough to melt boulders. Known to cause forest fires unintentionally. Comes with a free fire extinguisher and homeowner's insurance disclaimer.",
			Price: 14999, Image: "/pokemon/charizard.png", HP: 78, Attack: 84, Defense: 78, Speed: 100,
		},
		{
			ID: 9, Name: "Blastoise", Type: []string{"Water"},
			Description: "It crushes its foe under its heavy body to cause fainting. It is then free to do as it pleases. The water cannons on its shell are NOT for watering your garden. Okay, maybe a little.",
			Price: 12999, Image: "/pokemon/blastoise.png", HP: 79, Attack: 83, Defense: 100, Speed: 78,
		},
		{
			ID: 3, Name: "Venusaur", Type: []string{"Grass", "Poison"},
			Description: "The plant blooms when it absorbs solar energy. It stays on the move to seek sunlight. Makes an excellent conversation piece at dinner parties. 'Is that a flower or a dinosaur?' Yes.",
			Price: 11999, Image: "/pokemon/venusaur.png", HP: 80, Attack: 82, Defense: 83, Speed: 80,
		},
		{
			ID: 39, Name: "Jigglypuff", Type: []string{"Normal", "Fairy"},
			Description: "When its huge eyes waver, it sings a mysteriously soothing melody that lulls its enemies to sleep. Perfect for insomniacs, colicky babies, and boring meetings.",
			Price: 1999, Image: "/pokemon/jigglypuff.png", HP: 115, Attack: 45, Defense: 20, Speed: 20,
		},
		{
			ID: 143, Name: "Snorlax", Type: []string{"Normal"},
			Description: "Very lazy. Just eats and sleeps. As its rotund bulk builds, it becomes steadily more slothful. Basically your spirit animal on a Sunday afternoon. Endorsed by couch manufacturers worldwide.",
			Price: 9999, Image: "/pokemon/snorlax.png", HP: 160, Attack: 110, Defense: 65, Speed: 30,
		},
		{
			ID: 150, Name: "Mewtwo", Type: []string{"Psychic"},
			Description: "A Pokemon created by recombining Mew's genes. It's said to have the most savage heart among Pokemon. Perfect for winning arguments — it literally reads minds. No refunds.",
			Price: 49999, Image: "/pokemon/mewtwo.png", HP: 106, Attack: 110, Defense: 90, Speed: 130,
		},
		{
			ID: 151, Name: "Mew", Type: []string{"Psychic"},
			Description: "So rare that it is still said to be a mirage by many experts. Only a few people have seen it worldwide. If you're reading this, congratulations — you found one. Limited stock!",
			Price: 59999, Image: "/pokemon/mew.png", HP: 100, Attack: 100, Defense: 100, Speed: 100,
		},
		{
			ID: 133, Name: "Eevee", Type: []string{"Normal"},
			Description: "Its genetic code is irregular. It may mutate if exposed to radiation from element stones. The ultimate impulse buy — one purchase, eight possible evolutions. Gotta catch 'em all starts here.",
			Price: 6999, Image: "/pokemon/eevee.png", HP: 55, Attack: 55, Defense: 50, Speed: 55,
		},
		{
			ID: 94, Name: "Gengar", Type: []string{"Ghost", "Poison"},
			Description: "Under a full moon, this Pokemon likes to mimic the shadows of people and laugh at their fright. Great for Halloween, haunted houses, or scaring your roommate at 3 AM.",
			Price: 7999, Image: "/pokemon/gengar.png", HP: 60, Attack: 65, Defense: 60, Speed: 110,
		},
		{
			ID: 130, Name: "Gyarados", Type: []string{"Water", "Flying"},
			Description: "Rarely seen in the wild. Huge and vicious, it is capable of destroying entire cities in a rage. NOT suitable for apartment living. Requires a pool. A very large pool.",
			Price: 13999, Image: "/pokemon/gyarados.png", HP: 95, Attack: 125, Defense: 79, Speed: 81,
		},
		{
			ID: 149, Name: "Dragonite", Type: []string{"Dragon", "Flying"},
			Description: "It is said to make its home somewhere in the sea. It guides crews of shipwrecks to shore. Also delivers mail faster than any postal service. Basically a flying UPS truck with a smile.",
			Price: 19999, Image: "/pokemon/dragonite.png", HP: 91, Attack: 134, Defense: 95, Speed: 80,
		},
		{
			ID: 131, Name: "Lapras", Type: []string{"Water", "Ice"},
			Description: "A gentle soul that can read the minds of people. It ferries people across bodies of water. Uber Pool but make it Pokemon. Seats 4 comfortably, 6 if you're friendly.",
			Price: 11999, Image: "/pokemon/lapras.png", HP: 130, Attack: 85, Defense: 80, Speed: 60,
		},
		{
			ID: 54, Name: "Psyduck", Type: []string{"Water"},
			Description: "Always tormented by headaches. It uses its mysterious powers when its headache turns severe. Basically every developer during a production outage. Relatable and affordable.",
			Price: 1499, Image: "/pokemon/psyduck.png", HP: 50, Attack: 52, Defense: 48, Speed: 55,
		},
		{
			ID: 52, Name: "Meowth", Type: []string{"Normal"},
			Description: "Adores circular objects. Wanders the streets on a nightly basis to look for dropped loose change. Will literally pay for itself over time. Best ROI in the entire catalog.",
			Price: 2499, Image: "/pokemon/meowth.png", HP: 40, Attack: 45, Defense: 35, Speed: 90,
		},
		{
			ID: 74, Name: "Geodude", Type: []string{"Rock", "Ground"},
			Description: "Found in fields and mountains. Mistaking them for boulders, people often step or trip on them. The original pet rock, but this one punches back. Perfect paperweight.",
			Price: 999, Image: "/pokemon/geodude.png", HP: 40, Attack: 80, Defense: 100, Speed: 20,
		},
		{
			ID: 129, Name: "Magikarp", Type: []string{"Water"},
			Description: "Famous for being virtually useless. It can only splash and flop around. But hey, it's the cheapest Pokemon we have! Great gag gift. Evolves into Gyarados (sold separately for 139.99).",
			Price: 199, Image: "/pokemon/magikarp.png", HP: 20, Attack: 10, Defense: 55, Speed: 80,
		},
		{
			ID: 37, Name: "Vulpix", Type: []string{"Fire"},
			Description: "At the time of its birth, it has just one snow-white tail. The tail splits from the tip as it grows older. Instagram's most photogenic Pokemon. Guaranteed to get likes.",
			Price: 4499, Image: "/pokemon/vulpix.png", HP: 38, Attack: 41, Defense: 40, Speed: 65,
		},
		{
			ID: 26, Name: "Raichu", Type: []string{"Electric"},
			Description: "Its long tail serves as a ground to protect itself from its own high-voltage power. Pikachu's cooler older sibling that nobody talks about. Deserves more love.",
			Price: 10999, Image: "/pokemon/raichu.png", HP: 60, Attack: 90, Defense: 55, Speed: 110,
		},
		{
			ID: 144, Name: "Articuno", Type: []string{"Ice", "Flying"},
			Description: "A legendary bird Pokemon that is said to appear to doomed people who are lost in icy mountains. Also makes a fantastic air conditioner during summer. Very eco-friendly.",
			Price: 39999, Image: "/pokemon/articuno.png", HP: 90, Attack: 85, Defense: 100, Speed: 85,
		},
		{
			ID: 145, Name: "Zapdos", Type: []string{"Electric", "Flying"},
			Description: "A legendary bird Pokemon that is said to appear from clouds while dropping enormous lightning bolts. Powers a small city. Your electricity bill will thank you.",
			Price: 39999, Image: "/pokemon/zapdos.png", HP: 90, Attack: 90, Defense: 85, Speed: 100,
		},
		{
			ID: 146, Name: "Moltres", Type: []string{"Fire", "Flying"},
			Description: "Known as the legendary bird of fire. Every flap of its wings creates a dazzling flash of flames. Makes s'mores effortless. Camping will never be the same.",
			Price: 39999, Image: "/pokemon/moltres.png", HP: 90, Attack: 100, Defense: 90, Speed: 90,
		},
	}
}
