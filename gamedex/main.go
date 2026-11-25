package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gamedex/client"
)

var templates *template.Template

type pageAccueil struct {
	Title string
}

type pageJeux struct {
	Title string
	Query string
	Games []client.Game
}

type pageJeu struct {
	Title       string
	Game        client.Game
	Tournaments []client.Tournament
}

type pageErreur struct {
	Title   string
	Code    int
	Message string
}

func render(w http.ResponseWriter, status int, name string, data any) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		log.Printf("template %q : %v", name, err)
		http.Error(w, "erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func renderError(w http.ResponseWriter, status int, message string) {
	render(w, status, "erreur", pageErreur{
		Title:   fmt.Sprintf("Erreur %d", status),
		Code:    status,
		Message: message,
	})
}

func accueil(w http.ResponseWriter, r *http.Request) {
	render(w, http.StatusOK, "accueil", pageAccueil{Title: "Accueil"})
}

func listeJeux(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	games, err := client.FetchGames()
	if err != nil {
		log.Printf("FetchGames : %v", err)
		renderError(w, http.StatusServiceUnavailable,
			"Le service des jeux est indisponible pour le moment. Réessayez dans un instant.")
		return
	}

	render(w, http.StatusOK, "jeux", pageJeux{Title: "Les jeux", Query: q, Games: games})
}

func detailJeu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		renderError(w, http.StatusBadRequest, "L'identifiant d'un jeu doit être un nombre.")
		return
	}

	game, err := client.FetchGame(id)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			renderError(w, http.StatusNotFound, "Ce jeu n'existe pas.")
			return
		}
		log.Printf("FetchGame(%d) : %v", id, err)
		renderError(w, http.StatusServiceUnavailable,
			"Le service des jeux est indisponible pour le moment. Réessayez dans un instant.")
		return
	}

	tournaments, err := client.FetchTournaments(id)
	if err != nil {
		log.Printf("FetchTournaments(%d) : %v", id, err)
		renderError(w, http.StatusServiceUnavailable,
			"Les tournois de ce jeu sont indisponibles pour le moment. Réessayez dans un instant.")
		return
	}

	render(w, http.StatusOK, "jeu", pageJeu{Title: game.Name, Game: game, Tournaments: tournaments})
}

func introuvable(w http.ResponseWriter, r *http.Request) {
	renderError(w, http.StatusNotFound, "Cette page n'existe pas.")
}

func main() {
	port := flag.Int("port", 8080, "port d'écoute du serveur")
	flag.Parse()

	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("impossible de charger les templates (lancez le serveur depuis gamedex/) : %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /", introuvable)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("GameDex écoute sur http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
