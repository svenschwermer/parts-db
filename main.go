package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/svenschwermer/parts-db/auth"
	"github.com/svenschwermer/parts-db/config"
	"github.com/svenschwermer/parts-db/handler"
)

var templates = template.Must(template.ParseGlob("html/*.*"))

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

	auther := auth.New(templates)
	h, err := handler.New(templates, auther, db)
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

	err = http.ListenAndServe(config.Env.ListenAddress, h)
	log.Fatal(err)
}
