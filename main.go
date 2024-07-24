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

func main() {
	http.HandleFunc("/artists", artistsHandler)
	http.HandleFunc("/artist/", artistDetailHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server listening on port 8060...")
	fmt.Println("http://localhost:8060/artists")
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

	// Combine artist and relations data
	data := struct {
		Artist    Artists
		Relations Relations
	}{
		Artist:    artist,
		Relations: relations,
	}
	
	// Execute HTML template for artist details
	tmpl := template.Must(template.ParseFiles("artist_detail.html"))
	err = tmpl.Execute(w, data)
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
	tmpl := template.Must(template.ParseFiles("index.html")) // Adjust path as needed
	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
