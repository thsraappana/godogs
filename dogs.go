package main

import (
	"html/template"
	"log"
	"net/http"
)

type Page struct {
	Title     string
	BreedList []string
	ImageURL  string
}

var templates = template.Must(template.ParseFiles("tmpl/home.html"))

func home(w http.ResponseWriter, r *http.Request) {
	page := Page{"Dog app with Go!", nil, ""}

	err := templates.ExecuteTemplate(w, "home.html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
}

// TODO: fetch a list of breeds and set to BreedList
func handleFetchBreedsButtonClick(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle click event for breeds!")
}

// TODO: fetch a random image of a dog and set to imageURL
func handleFetchRandomImageButtonClick(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle click event for random image!")

	page := Page{
		Title:     "Dog app with Go!",
		BreedList: nil,
		ImageURL:  "woof",
	}
	err := templates.ExecuteTemplate(w, "home.html", page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/handle-fetch-breeds-button-click", handleFetchBreedsButtonClick)
	http.HandleFunc("/handle-random-image-button-click", handleFetchRandomImageButtonClick)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
