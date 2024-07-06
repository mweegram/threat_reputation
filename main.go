package main

import (
	"net/http"

	"github.com/mweegram/threat_reputation/handlers"
)

func main() {

	http.HandleFunc("/", handlers.HomepageHandler)
	http.ListenAndServe("0.0.0.0:5000", nil)
}
