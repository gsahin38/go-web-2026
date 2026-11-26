package catalog

type Game struct {
	ID    int
	Name  string
	Genre string
	Year  int
}

type Tournament struct {
	ID     int
	Name   string
	GameID int
	City   string
	Date   string
}

type Catalog struct {
	games       []Game
	gamesByID   map[int]Game
	tournaments map[int][]Tournament // clé : ID du jeu
}

func New(games []Game, tournaments map[int][]Tournament) *Catalog {
	byID := make(map[int]Game, len(games))
	for _, g := range games {
		byID[g.ID] = g
	}
	return &Catalog{games: games, gamesByID: byID, tournaments: tournaments}
}

func (c *Catalog) Games() []Game {
	return c.games
}

func (c *Catalog) Game(id int) (Game, bool) {
	g, ok := c.gamesByID[id]
	return g, ok
}

func (c *Catalog) TournamentsOf(gameID int) []Tournament {
	return c.tournaments[gameID]
}
