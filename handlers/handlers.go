package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/mweegram/threat_reputation/logic"
)

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./templates/homepage.html"))
	stats := logic.HomepageStats()
	temp.Execute(w, stats)
}

func NewThreatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Bad Method"))
		return
	}

	temp := template.Must(template.ParseFiles("./templates/notification.html"))
	new_threat := logic.Threat{
		Filename:  r.FormValue("filename"),
		Sha256:    r.FormValue("filehash"),
		Submitted: time.Now().Format("2006-01-02"),
	}

	err := new_threat.Validate()
	if err != nil {
		data := logic.Notification{
			Status:  "Error",
			Message: fmt.Sprintf("%v", err),
			Colour:  "red",
		}
		temp.Execute(w, data)
		return
	}

	err = logic.DB_INSTANCE.CreateThreat(new_threat)
	if err != nil {
		data := logic.Notification{
			Status:  "Error",
			Message: fmt.Sprintf("%v", err),
			Colour:  "red",
		}
		temp.Execute(w, data)
		return
	}

	//TODO: Add Duplication Checking - Simple, But What's The Best Way?

	data := logic.Notification{
		Status:  "Success",
		Message: "New Threat Added to Database",
		Colour:  "green",
	}
	temp.Execute(w, data)
}
