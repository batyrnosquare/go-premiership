package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

type APIHandler struct {
	Template *template.Template
}

func NewAPIHandler(templatePath string) (*APIHandler, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &APIHandler{
		Template: tmpl,
	}, nil
}

func (h *APIHandler) HandleAPI(w http.ResponseWriter, r *http.Request) {
	url := "https://api.sportmonks.com/v3/football/fixtures/18535517"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Add("Authorization", "PWLlkCQssouNBTLciwRmJDhI2xcwFTZPoKsIH820H1JI0wPix4kas9JeH4MP")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		http.Error(w, readErr.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Render the API response using the template
	err = h.Template.Execute(w, string(body))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
