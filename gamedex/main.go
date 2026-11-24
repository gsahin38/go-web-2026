package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Jeu struct {
	ID     int
	Nom    string
	Annee  int
	Studio string
	Genres []string
}

var jeux = []Jeu{
	{1, "Celeste", 2018, "Maddy Makes Games", []string{"plateforme", "indé"}},
	{2, "Hades", 2020, "Supergiant Games", []string{"rogue-like", "action"}},
	{3, "Rocket League", 2015, "Psyonix", []string{"sport", "course"}},
	{4, "Street Fighter 6", 2023, "Capcom", []string{"combat"}},
	{5, "Stardew Valley", 2016, "ConcernedApe", []string{"simulation", "indé"}},
	{6, "League of Legends", 2009, "Riot Games", []string{"moba", "stratégie"}},
	{7, "Valorant", 2020, "Riot Games", []string{"fps", "tactique"}},
	{8, "Dead Cells", 2018, "Motion Twin", []string{"rogue-like", "action", "indé"}},
	{9, "Trackmania", 2020, "Nadeo", []string{"course"}},
	{10, "Mario Kart 8 Deluxe", 2017, "Nintendo", []string{"course", "party"}},
	{11, "Hollow Knight", 2017, "Team Cherry", []string{"metroidvania", "action", "indé"}},
	{12, "Minecraft", 2011, "Mojang", []string{"sandbox", "survie"}},
}

const accueil = `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="utf-8">
  <title>GameDex</title>
  <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
  <header><h1>GameDex</h1></header>
  <main>
    <p>Le catalogue de jeux vidéo.</p>
    <p><a class="bouton" href="/jeux">Voir les jeux</a></p>
  </main>
</body>
</html>`

func handleAccueil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, accueil)
}

func handleJeux(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, j := range jeux {
		fmt.Fprintf(w, "%d. %s (%d) - %s\n", j.ID, j.Nom, j.Annee, j.Studio)
	}
}

func handleJeu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "400 : l'identifiant doit être un nombre", http.StatusBadRequest)
		return
	}
	for _, j := range jeux {
		if j.ID == id {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			fmt.Fprintf(w, "%s\nAnnée : %d\nStudio : %s\nGenres : %s\n",
				j.Nom, j.Annee, j.Studio, strings.Join(j.Genres, ", "))
			return
		}
	}
	http.Error(w, "404 : jeu introuvable", http.StatusNotFound)
}

func main() {
	port := flag.Int("port", 8080, "port d'écoute du serveur")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleAccueil)
	mux.HandleFunc("GET /jeux", handleJeux)
	mux.HandleFunc("GET /jeux/{id}", handleJeu)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("GameDex sur http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}