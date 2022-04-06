package web

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"self-update/updater"
	"text/template"
)

//go:embed templates
var indexHTML embed.FS

type Server struct {
	Port    int
	Host    string
	Updater updater.Updater
}

func (s *Server) backgroundUpdate() error {
	err := s.Updater.Update()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return nil
}

func (s *Server) Serve() error {
	status := Status{
		Version: s.Updater.CurrentVersion,
	}
	log.Printf("App Version %s \n", status.Version)

	// Load page from embedded template
	page, err := template.ParseFS(indexHTML, "templates/index.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		if err := page.Execute(w, status); err != nil {
			w.WriteHeader(500)
			log.Fatal(err)
		}
	})

	mux.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		err := s.Updater.SetLatestReleaseInfo()
		if err != nil {
			log.Fatal(err)
		}

		if s.Updater.NeedUpdate() {
			status.NewVersion = s.Updater.LatestRelease.TagName
		}
		if err := page.Execute(w, status); err != nil {
			w.WriteHeader(500)
			log.Fatal(err)
		}
	})

	mux.HandleFunc("/install", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusOK)
		go s.backgroundUpdate()
	})

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	log.Printf("Listening on %s...\n", addr)

	return http.ListenAndServe(addr, mux)
}
