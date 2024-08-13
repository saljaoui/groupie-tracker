package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

var tmpl *template.Template

func main() {
	var err error
	tmpl, err = template.ParseGlob("templates/html/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", artistsHandler)
	http.HandleFunc("/artist/", artistDetailHandler)

	fmt.Println("Server listening on port 8060...")
	fmt.Println("http://localhost:8060")
	log.Fatal(http.ListenAndServe(":8060", nil))
}

func artistDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/artist/"):]

	// Fetch artist details
	artistResp, err := http.Get(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%s", id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer artistResp.Body.Close()

	var artist Artists
	if err := json.NewDecoder(artistResp.Body).Decode(&artist); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch relations
	relationsResp, err := http.Get(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%s", id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer relationsResp.Body.Close()

	var relations Relations
	if err := json.NewDecoder(relationsResp.Body).Decode(&relations); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute HTML template for artist details
	err = tmpl.ExecuteTemplate(w, "artist_detail.html", artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func artistsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var artists []Artists
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute HTML template
	err = tmpl.ExecuteTemplate(w, "index.html", artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}