package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	FetchImages = "fetchImg"
	FetchBreeds = "fetchBreeds"
)

type Page struct {
	Title     string
	BreedList map[string][]string
	ImageURL  string
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func fetch(endpoint string) ([]byte, error) {
	// create an HTTP client
	client := &http.Client{}

	// create a GET request
	req, err := http.NewRequest("GET", "https://dog.ceo/api/breeds"+endpoint, nil)
	if err != nil {
		panic(err)
	}

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body, err
}

func (p *Page) fetchImage() (string, error) {
	body, err := fetch("/image/random")

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	return data.Message, err
}

func (p *Page) fetchList() (map[string][]string, error) {
	body, _ := fetch("/list/all")

	breedMap := make(map[string][]string)

	var data json.RawMessage
	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	message, ok := response["message"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to parse message")
	}

	messageMap := make(map[string]string)
	for k, v := range message {
		messageMap[k] = fmt.Sprintf("%v", v)
	}

	for key, value := range messageMap {
		if value != "[]" {
			s := value
			s = strings.Trim(s, "[]") // remove brackets from the string
			arr := strings.Split(s, " ")
			breedMap[key] = arr
		} else {
			breedMap[key] = nil
		}
	}

	return breedMap, nil
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

func eventHandler(eventType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use myParam in the handler function
		page := Page{
			Title:     "Dog app with Go!",
			BreedList: nil,
			ImageURL:  "",
		}

		switch eventType {
		case FetchBreeds:
			breedMap, _ := page.fetchList()
			page.BreedList = breedMap
		case FetchImages:
			img, _ := page.fetchImage()
			page.ImageURL = img
		}

		err := templates.ExecuteTemplate(w, "home.html", page)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", home)
	http.HandleFunc("/handle-fetch-breeds-button-click", eventHandler(FetchBreeds))
	http.HandleFunc("/handle-random-image-button-click", eventHandler(FetchImages))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
