package main

import (
	"fmt"
	"log"
	
	"gamedex/client"
)

func main() {
	log.SetFlags(0) 

	games, err := client.FetchGames()
	if err != nil {
		log.Fatalf("Erreur : %v", err) 
	}

	for _, g := range games {
		tournaments, err := client.FetchTournaments(g.ID)
		if err != nil {
			log.Fatalf("Erreur : %v", err)
		}
		fmt.Printf("%s (%d) : %d tournoi(s)\n", g.Name, g.Year, len(tournaments))
	}
}
