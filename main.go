package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
	_ "github.com/mattn/go-sqlite3"

	"github.com/svenschwermer/parts-db/auth"
	"github.com/svenschwermer/parts-db/config"
	"github.com/svenschwermer/parts-db/server"
)

var templates = template.Must(template.ParseGlob("html/*.*"))

func watchTemplates(watcher *fsnotify.Watcher) {
	for e := range watcher.Events {
		log.Println("Template file event:", e)
		tmpl, err := template.ParseGlob("html/*.*")
		if err != nil {
			log.Println("Failed to parse templates:", err)
		} else {
			*templates = *tmpl
		}
	}
}

func main() {
	config.Process()

	db, err := sql.Open("sqlite3", config.Env.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		log.Fatal(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go watchTemplates(watcher)
	err = watcher.Add("html")
	if err != nil {
		log.Fatal(err)
	}

	auther := auth.New(templates)
	h, err := server.New(templates, auther, db)
	if err != nil {
		log.Fatal(err)
	}

	h.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("html/assets"))))
	h.HandleFunc(auth.Path, auther.Login)
	h.HandleFunc("/list", h.List)
	h.HandleFunc("/change-inventory", h.ChangeInventory)
	h.HandleFunc("/new", h.New)
	h.HandleFunc("/edit", h.Edit)
	h.HandleFunc("/mouser/", h.Mouser)
	h.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			http.Redirect(w, req, "/list", http.StatusFound)
		} else {
			http.NotFound(w, req)
		}
	})

	err = http.ListenAndServe(config.Env.ListenAddress, h)
	log.Fatal(err)
}
