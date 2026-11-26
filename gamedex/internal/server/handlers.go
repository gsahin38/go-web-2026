package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gamedex/internal/catalog"
)

type pageAccueil struct {
	Title string
}

type pageJeux struct {
	Title string
	Query string 
	Games []catalog.Game
}

type pageJeu struct {
	Title       string
	Game        catalog.Game
	Tournaments []catalog.Tournament
}

type pageErreur struct {
	Title   string
	Code    int
	Message string
}

func (s *Server) render(w http.ResponseWriter, status int, name string, data any) {
	var buf bytes.Buffer
	if err := s.templates.ExecuteTemplate(&buf, name, data); err != nil {
		log.Printf("template %q : %v", name, err)
		http.Error(w, "erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status) 
	buf.WriteTo(w)
}

func (s *Server) renderError(w http.ResponseWriter, status int, message string) {
	s.render(w, status, "erreur", pageErreur{
		Title:   fmt.Sprintf("Erreur %d", status),
		Code:    status,
		Message: message,
	})
}

func (s *Server) accueil(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusOK, "accueil", pageAccueil{Title: "Accueil"})
}

func (s *Server) listeJeux(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusOK, "jeux", pageJeux{
		Title: "Les jeux",
		Query: r.URL.Query().Get("q"),
		Games: s.cat.Games(),
	})
}

func (s *Server) detailJeu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		s.renderError(w, http.StatusBadRequest, "L'identifiant d'un jeu doit être un nombre.")
		return
	}

	game, ok := s.cat.Game(id)
	if !ok {
		s.renderError(w, http.StatusNotFound, "Ce jeu n'existe pas.")
		return
	}

	s.render(w, http.StatusOK, "jeu", pageJeu{
		Title:       game.Name,
		Game:        game,
		Tournaments: s.cat.TournamentsOf(id),
	})
}

func (s *Server) introuvable(w http.ResponseWriter, r *http.Request) {
	s.renderError(w, http.StatusNotFound, "Cette page n'existe pas.")
}
