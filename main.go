package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"

	"github.com/svenschwermer/parts-db/auth"
	"github.com/svenschwermer/parts-db/handler"
)

var (
	templates = template.Must(template.ParseGlob("html/*.*"))
	env       = struct {
		ListenAddress string `envconfig:"LISTEN_ADDRESS" default:":80"`
		SitePassword  string `envconfig:"SITE_PASSWORD" required:"true"`
		DatabasePath  string `envconfig:"DB_PATH" default:"parts.db"`
	}{}
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		envconfig.Usage("", &env)
		return
	}
	envconfig.MustProcess("", &env)

	db, err := sql.Open("sqlite3", env.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	auther := auth.New(templates, env.SitePassword)
	h := handler.New(templates, auther)

	h.HandleFunc(auth.Path, auther.Login)
	h.HandleFunc(handler.ListPath, h.List)
	h.HandleFunc(handler.ChangeInventoryPath, h.ChangeInventory)

	err = http.ListenAndServe(env.ListenAddress, h)
	log.Fatal(err)
}
