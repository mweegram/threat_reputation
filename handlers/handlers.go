package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/mweegram/threat_reputation/logic"
)

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./templates/homepage.html"))
	content := logic.Homepage_Content{
		Stats:   logic.HomepageStats(),
		Threats: logic.Recent_Threats(),
	}
	temp.Execute(w, content)
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

func GetThreatHandler(w http.ResponseWriter, r *http.Request) {
	threat_id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		temp := template.Must(template.ParseFiles("./templates/notification.html"))
		data := logic.Notification{
			Status:  "Error",
			Message: "Invalid ID Parameter Provided",
			Colour:  "red",
		}
		temp.Execute(w, data)
		return
	}

	enriched_comments := logic.Get_Comments(threat_id)

	selected_threat, err := logic.DB_INSTANCE.GetThreat(threat_id)
	if err != nil {
		temp := template.Must(template.ParseFiles("./templates/notification.html"))
		data := logic.Notification{
			Status:  "Error",
			Message: fmt.Sprintf("Error: %v", err),
			Colour:  "red",
		}
		temp.Execute(w, data)
		return
	}

	temp := template.Must(template.ParseFiles("./templates/threat.html"))
	data := logic.DisplayThreat{
		ID:        selected_threat.ID,
		Filename:  selected_threat.Filename,
		Sha256:    selected_threat.Sha256,
		Comments:  enriched_comments,
		Submitted: selected_threat.Submitted,
	}
	temp.Execute(w, data)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Bad Method"))
		return
	}

	comment_value := r.FormValue("comment") //TODO: Validate Length
	threat_id, _ := strconv.Atoi(r.PathValue("id"))

	err := logic.Add_Comment(threat_id, comment_value)
	if err != nil {
		temp := template.Must(template.ParseFiles("./templates/notification.html"))
		data := logic.Notification{
			Status:  "Error",
			Message: fmt.Sprintf("Error: %v", err),
			Colour:  "red",
		}
		temp.Execute(w, data)
		return
	}

	temp := template.Must(template.ParseFiles("./templates/comments.html"))
	data := logic.Get_Comments(threat_id)
	temp.Execute(w, data)
}
