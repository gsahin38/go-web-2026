package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8081"

const gamesPath = "/api/games"

const timeout = 3 * time.Second

var httpClient = &http.Client{Timeout: timeout}

var ErrNotFound = errors.New("ressource introuvable")

type Game struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Genre string `json:"genre"`
	Year  int    `json:"year"`
}

type Tournament struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	GameID int    `json:"game_id"`
	City   string `json:"city"`
	Date   string `json:"date"`
}

func FetchGames() ([]Game, error) {
	var games []Game
	if err := getJSON(gamesPath, &games); err != nil {
		return nil, err
	}
	for _, g := range games {
		if err := checkGame(g); err != nil {
			return nil, err
		}
	}
	return games, nil
}

func FetchGame(id int) (Game, error) {
	var game Game
	if err := getJSON(fmt.Sprintf("/api/games/%d", id), &game); err != nil {
		return Game{}, err
	}
	if err := checkGame(game); err != nil {
		return Game{}, err
	}
	return game, nil
}

func FetchTournaments(gameID int) ([]Tournament, error) {
	var tournaments []Tournament
	path := fmt.Sprintf("/api/games/%d/tournaments", gameID)
	if err := getJSON(path, &tournaments); err != nil {
		return nil, err
	}
	return tournaments, nil
}

func checkGame(g Game) error {
	if g.Name == "" {
		return fmt.Errorf("le jeu n°%d n'a pas de nom : les tags json de Game ne correspondent pas au JSON de l'API", g.ID)
	}
	return nil
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

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%s : %w", path, ErrNotFound)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("l'API a répondu %s pour %s", resp.Status, path)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("le JSON reçu pour %s n'a pas la forme attendue : %w", path, err)
	}
	return nil
}
