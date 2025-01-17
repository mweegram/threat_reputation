package main

import (
	"net/http"

	"github.com/mweegram/threat_reputation/database"
	"github.com/mweegram/threat_reputation/handlers"
	"github.com/mweegram/threat_reputation/logic"
)

func main() {
	logic.DB_INSTANCE.DB = database.DatabaseConnect()
	http.HandleFunc("/", handlers.HomepageHandler)
	http.HandleFunc("/newthreat", handlers.NewThreatHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/threat/{id}", handlers.GetThreatHandler)
	http.HandleFunc("/threat/{id}/comment", handlers.CreateCommentHandler)
	http.ListenAndServe("0.0.0.0:5000", nil)
}
