package server

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"gamedex/internal/catalog"
)

type Server struct {
	cat       *catalog.Catalog
	templates *template.Template
	mux       *http.ServeMux
}

func New(cat *catalog.Catalog, webDir string) (*Server, error) {
	templates, err := template.ParseGlob(filepath.Join(webDir, "templates", "*.html"))
	if err != nil {
		return nil, fmt.Errorf("chargement des templates : %w", err)
	}

	s := &Server{cat: cat, templates: templates, mux: http.NewServeMux()}

	s.mux.HandleFunc("GET /{$}", s.accueil)
	s.mux.HandleFunc("GET /jeux", s.listeJeux)
	s.mux.HandleFunc("GET /jeux/{id}", s.detailJeu)

	// /static/css/style.css -> web/static/css/style.css
	staticDir := http.Dir(filepath.Join(webDir, "static"))
	s.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	s.mux.HandleFunc("GET /", s.introuvable)

	return s, nil
}

func (s *Server) Handler() http.Handler {
	return s.mux
}
