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
	FetchRandImg  = "fetchRandImg"
	FetchBreeds   = "fetchBreeds"
	FetchBreedImg = " fetchBreedImg"
)

type Page struct {
	Title        string
	BreedList    map[string][]string
	ImageURL     string
	Error        bool
	ErrorMessage string
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func fetch(endpoint string) ([]byte, error) {
	// create an HTTP client
	client := &http.Client{}

	// create a GET request
	req, err := http.NewRequest("GET", "https://dog.ceo/api"+endpoint, nil)
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

func (p *Page) fetchImage(breed string) (string, error) {
	var endpoint string
	if len(breed) == 0 {
		endpoint = "/breeds/image/random"
	} else {
		endpoint = "/breed/" + breed + "/images/random"
	}
	body, err := fetch(endpoint)
	if err != nil {
		fmt.Println("error in image fetch")
		return "", err
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	return data.Message, nil
}

func (p *Page) fetchList() (map[string][]string, error) {
	body, _ := fetch("/breeds/list/all")

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
	page := Page{"Dog app with Go!", nil, "", false, ""}

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
		case FetchRandImg:
			img, _ := page.fetchImage("")
			page.ImageURL = img
		case FetchBreedImg:
			// read the body of the request
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			// unmarshal the body data to read the values
			var breedData map[string]interface{}
			err = json.Unmarshal([]byte(string(body)), &breedData)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			data, err := page.fetchImage(breedData["breed"].(string))
			if data == "Breed not found (master breed does not exist)" {
				page.Error = true
				page.ErrorMessage = data
			} else {
				page.ImageURL = data
			}
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
	http.HandleFunc("/fetch-breeds", eventHandler(FetchBreeds))
	http.HandleFunc("/fetch-random-image", eventHandler(FetchRandImg))
	http.HandleFunc("/fetch-breed-image", eventHandler(FetchBreedImg))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
