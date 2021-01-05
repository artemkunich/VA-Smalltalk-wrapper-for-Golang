package main

import "net/http"

type Request struct {
	HttpRequest *http.Request
	Resource    string
	Method      string
	Flags       int //0b01-pass request content as argument to Smalltalk method, 0b10-in Smalltalk should be created class instance
}

type Response struct {
	Status  int
	Content string
}

var RequestChan chan Request
var ResponseChan chan Response

func InitChannels() {
	RequestChan = make(chan Request)
	ResponseChan = make(chan Response)
}
