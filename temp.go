package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type ViewData struct {
	Title   string `json:"Title"`
	Message string `json:"Message"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read data from a JSON file (you should replace 'data.json' with your actual JSON file)
		jsonFile, err := os.Open("laptops.json")
		if err != nil {
			http.Error(w, "Error reading JSON file", http.StatusInternalServerError)
			return
		}
		defer jsonFile.Close()

		var data []ViewData
		decoder := json.NewDecoder(jsonFile)
		err = decoder.Decode(&data)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
			return
		}

		// Create a template that renders each item in the list
		tmpl, err := template.ParseFiles("products.html")

		// Execute the template with the list of ViewData items
		tmpl.Execute(w, data)
	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
