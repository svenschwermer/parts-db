package digikey

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func Main() {
	code := make(chan string)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("state=%s code=%s", r.FormValue("state"), r.FormValue("code"))
		code <- r.FormValue("code")
	})
	go http.ListenAndServe(":8080", nil)

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "QNBLSjgo5EinV7kUogbU3KJKLciBYYok",
		ClientSecret: "GPSWOb2reDwnimfY",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://sandbox-api.digikey.com/v1/oauth2/authorize",
			TokenURL: "https://sandbox-api.digikey.com/v1/oauth2/token",
		},
		RedirectURL: "http://localhost:8080/callback",
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	tok, err := conf.Exchange(ctx, <-code)
	if err != nil {
		log.Fatal(err)
	}
	client := conf.Client(ctx, tok)
	req, err := http.NewRequest("GET", "https://sandbox-api.digikey.com/Search/v3/Products/P5555-ND", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-DIGIKEY-Client-Id", conf.ClientID)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}
