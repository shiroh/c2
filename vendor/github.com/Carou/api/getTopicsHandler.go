package api

import (
	"net/http"

	"github.com/Carou/models"
	"github.com/alecthomas/template"
)

// getTopicsHandler is the handler return the topics based by votes ascending.
func getTopicsHandler(w http.ResponseWriter, r *http.Request) {
	topics, err := cacheV1.GetTopics(20)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	renderTemplate(w, "index", topics)
	w.WriteHeader(200)
	return
}

func renderTemplate(w http.ResponseWriter, tmpl string, topics []models.Topic) {
	t, _ := template.ParseFiles("views/" + tmpl + ".html")
	t.Execute(w, topics)
}
