package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"gamedex/internal/catalog"
)

const baseURL = "http://localhost:8081"

const gamesPath = "/api/games"

const timeout = 3 * time.Second

var httpClient = &http.Client{Timeout: timeout}

type gameJSON struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Genre string `json:"genre"`
	Year  int    `json:"year"`
}

type tournamentJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
	Date string `json:"date"`
}

func FetchGames() ([]catalog.Game, error) {
	var raw []gameJSON
	if err := getJSON(gamesPath, &raw); err != nil {
		return nil, err
	}

	games := make([]catalog.Game, 0, len(raw))
	for _, g := range raw {
		if g.Name == "" {
			return nil, fmt.Errorf("le jeu n°%d n'a pas de nom : les tags json de gameJSON ne correspondent pas au JSON de l'API", g.ID)
		}
		games = append(games, catalog.Game{ID: g.ID, Name: g.Name, Genre: g.Genre, Year: g.Year})
	}
	return games, nil
}

func FetchTournaments(gameID int) ([]catalog.Tournament, error) {
	var raw []tournamentJSON
	path := fmt.Sprintf("/api/games/%d/tournaments", gameID)
	if err := getJSON(path, &raw); err != nil {
		return nil, err
	}

	tournaments := make([]catalog.Tournament, 0, len(raw))
	for _, t := range raw {
		tournaments = append(tournaments, catalog.Tournament{
			ID: t.ID, Name: t.Name, GameID: gameID, City: t.City, Date: t.Date,
		})
	}
	return tournaments, nil
}

func getJSON(path string, dest any) error {
	resp, err := httpClient.Get(baseURL + path)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return fmt.Errorf("l'API n'a pas répondu à %s en moins de %s (timeout) : %w", path, timeout, err)
		}
		return fmt.Errorf("impossible de joindre l'API pour %s (est-elle lancée ?) : %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("l'API a répondu %s pour %s", resp.Status, path)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("le JSON reçu pour %s n'a pas la forme attendue : %w", path, err)
	}
	return nil
}
