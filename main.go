package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jcuga/golongpoll"
)

func main() {
	manager, err := golongpoll.StartLongpoll(golongpoll.Options{
		LoggingEnabled:                 false,
		MaxLongpollTimeoutSeconds:      120,
		MaxEventBufferSize:             100,
		EventTimeToLiveSeconds:         60 * 2, // Event's stick around for 2 minutes
		DeleteEventAfterFirstRetrieval: false,
	})
	if err != nil {
		log.Fatalf("Failed to create manager: %q", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", AdvancedExampleHomepage)
	r.HandleFunc("/actions", getUserActionHandler(manager))
	r.HandleFunc("/events", getEvent(manager))
	fmt.Println("Serving webpage at http://127.0.0.1:8080/")
	http.ListenAndServe("127.0.0.1:8080", r)
	manager.Shutdown()
}

func AdvancedExampleHomepage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, nil)
}

func getUserActionHandler(manager *golongpoll.LongpollManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query().Get("text")
		if len(text) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error!"))
		}
		manager.Publish("mes", map[string]string{"text": text})
	}
}

func getEvent(manager *golongpoll.LongpollManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		manager.SubscriptionHandler(w, r)
	}
}
