package main

import (
	"encoding/json"
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
		EventTimeToLiveSeconds:         172800, // Event's stick around for 2 minutes
		DeleteEventAfterFirstRetrieval: false,
	})
	if err != nil {
		log.Fatalf("Failed to create manager: %q", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", AdvancedExampleHomepage)
	r.HandleFunc("/messages", getUserActionHandler(manager))
	r.HandleFunc("/events", getEvent(manager))
	r.Headers("Access-Control-Allow-Origin", "*")
	fmt.Println("Serving webpage at http://127.0.0.1:5500/")
	http.ListenAndServe("127.0.0.1:5500", r)
	manager.Shutdown()
}

func AdvancedExampleHomepage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, nil)
}

type test_st struct {
	Text string
}

func getUserActionHandler(manager *golongpoll.LongpollManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		decoder := json.NewDecoder(r.Body)
		var t test_st
		err := decoder.Decode(&t)
		text := t.Text
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(t.Text)
		if len(text) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error!"))
			return
		}
		manager.Publish("mes", map[string]string{"text": text})
	}
}

func getEvent(manager *golongpoll.LongpollManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		manager.SubscriptionHandler(w, r)
	}
}
