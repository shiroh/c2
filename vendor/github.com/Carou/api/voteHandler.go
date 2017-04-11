package api

import (
	"net/http"

	"github.com/Carou/models"
)

// voteHandler is the handler to process upvote and downvote
func voteHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	id := r.FormValue("id")
	if body == "1" {
		topic := models.Topic{ID: id}
		cacheV1.UpvoteTopic(topic)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if body == "-1" {
		topic := models.Topic{ID: id}
		cacheV1.DownvoteTopic(topic)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		w.WriteHeader(400)
		return
	}
}
