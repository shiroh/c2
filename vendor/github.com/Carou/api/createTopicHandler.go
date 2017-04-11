package api

import (
	"net/http"

	"github.com/Carou/models"
)

// createTopicHandler is the handler to create a new topic. Redirect to home page after finish processing
func createTopicHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	topic := models.NewTopic(body)
	storageV1.CreateTopic(topic)
	http.Redirect(w, r, "/", http.StatusFound)
	return
}
