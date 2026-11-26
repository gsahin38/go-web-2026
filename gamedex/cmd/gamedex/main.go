package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
	"gamedex/internal/server"
)

const webDir = "web"

func main() {
	port := flag.Int("port", 8080, "port d'écoute du serveur")
	flag.Parse()

	games, err := api.FetchGames()
	if err != nil {
		log.Fatalf("chargement des jeux impossible : %v", err)
	}

	tournaments := make(map[int][]catalog.Tournament, len(games))
	for _, g := range games {
		t, err := api.FetchTournaments(g.ID)
		if err != nil {
			log.Fatalf("chargement des tournois du jeu %d impossible : %v", g.ID, err)
		}
		tournaments[g.ID] = t
	}

	cat := catalog.New(games, tournaments)

	srv, err := server.New(cat, webDir)
	if err != nil {
		log.Fatalf("démarrage du serveur impossible : %v", err)
	}

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("GameDex : %d jeux chargés, écoute sur http://localhost%s", len(games), addr)
	log.Fatal(http.ListenAndServe(addr, srv.Handler()))
}
