package handlers

import (
	"html/template"
	"net/http"
)

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./templates/homepage.html"))
	temp.Execute(w, nil)
}
