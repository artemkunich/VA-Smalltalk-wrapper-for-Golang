package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func GoIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wellcome to the Go HTTP Server")
}

func GoTime(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go time: %v", time.Now())
}

func SmalltalkIndex(w http.ResponseWriter, r *http.Request) {
	RequestChan <- Request{HttpRequest: r, Resource: "GoHttpServer", Method: "index", Flags: 0b0}
	response := <-ResponseChan

	if response.Content != "" {
		fmt.Fprintf(w, response.Content)
	}

	if response.Status != 200 {
		w.WriteHeader(response.Status)
	}
}

func SmalltalkTime(w http.ResponseWriter, r *http.Request) {
	RequestChan <- Request{HttpRequest: r, Resource: "AbtTimestamp", Method: "now", Flags: 0b0}
	response := <-ResponseChan

	if response.Content != "" {
		fmt.Fprintf(w, "Smalltalk time: %v", response.Content)
	}

	if response.Status != 200 {
		w.WriteHeader(response.Status)
	}
}

var routes = Routes{
	Route{
		"GoIndex",
		"GET",
		"/Go/",
		GoIndex,
	},
	Route{
		"GoTime",
		"GET",
		"/Go/time",
		GoTime,
	},

	Route{
		"SmalltalkIndex",
		"GET",
		"/Smalltalk/",
		SmalltalkIndex,
	},
	Route{
		"SmalltalkTime",
		"GET",
		"/Smalltalk/time",
		SmalltalkTime,
	},
}
