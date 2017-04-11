package api

import (
	"net/http"
	"os"

	"github.com/Carou/cache"
	"github.com/Carou/storage"
	"github.com/Carou/utils"
	"github.com/gorilla/mux"
)

var (
	storageV1 = storage.NewV1()
	cacheV1   = cache.NewV1(storageV1)
)

// StartHttpServer start the http server.
var StartHttpServer = func() {
	r := mux.NewRouter()

	r.HandleFunc("/", getTopicsHandler).Methods("GET")
	r.HandleFunc("/create", createTopicHandler).Methods("POST")
	r.HandleFunc("/upvote", voteHandler).Methods("POST")
	r.HandleFunc("/downvote", voteHandler).Methods("POST")
	http.ListenAndServe(resolveAddress(), r)
}

var resolveAddress = func() string {
	if port := os.Getenv("PORT"); len(port) > 0 {
		utils.Info.Printf("Environment variable PORT=\"%s\"", port)
		return ":" + port
	}
	utils.Info.Println("Environment variable PORT is undefined. Using port :8080 by default")
	return ":8080"
}
