package main

import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"unsafe"

	"github.com/gorilla/mux"
)

var logFd *os.File

//export StartLogging
func StartLogging(f *C.char) int {
	var fs string = C.GoString(f)
	var err error

	if logFd, err = os.OpenFile(fs, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
		return -1
	}
	log.SetOutput(logFd)
	return 0
}

//export StopLogging
func StopLogging() int {
	if logFd == nil {
		return 0
	}
	if err := logFd.Close(); err != nil {
		return -1
	}
	logFd = nil
	return 0
}

//export RunGoHttpServer
func RunGoHttpServer(port int) int {
	InitChannels()
	router := NewRouter()
	go runGoHttpServer(port, router)

	return 0
}

func runGoHttpServer(port int, router *mux.Router) {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

//export PutResponse
func PutResponse(status C.int, res *C.char) int {
	var response string = C.GoString(res)
	ResponseChan <- Response{int(status), response}
	return 0
}

//export GetRequest
func GetRequest(resource *C.char, method *C.char, content *C.char, contentLength C.int) int {

	request := <-RequestChan

	if writeToBuffer(request.Resource, resource, 100) < 0 {
		ResponseChan <- Response{500, "Error occured 1"}
		return -1
	}
	if writeToBuffer(request.Method, method, 100) < 0 {
		ResponseChan <- Response{500, "Error occured 2"}
		return -2
	}

	if (request.Flags & 0b1) == 0b1 {

		body, err := ioutil.ReadAll(request.HttpRequest.Body)
		if err != nil {
			log.Printf("%v Error: %v", time.Now(), err.Error())
			ResponseChan <- Response{500, "Error occured 3"}
			return -3
		}

		if writeToBuffer(string(body), content, int(contentLength)) < 0 {
			ResponseChan <- Response{500, "Error occured 4"}
			return -4
		}
		return request.Flags
	}

	return request.Flags
}

func writeToBuffer(content string, buffer *C.char, bufferSize int) int {
	if bufferSize < len(content)+1 {
		return -1
	}

	i := 0
	for ; i < len(content); i++ {
		*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(buffer)) + uintptr(i))) = C.char(content[i])
	}
	*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(buffer)) + uintptr(i))) = C.char(0) //null terminator
	return 0
}

func main() {}
