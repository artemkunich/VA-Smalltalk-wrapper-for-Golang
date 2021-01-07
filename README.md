# Common

This repository contains experimental Cgo library and Smalltalk wrapper, which show how to use VA Smalltalk functionality from a Go HTTP server (or other Go code). 

# How it works

Go runtime and Smalltalk virtual machine run in the same process. Smalltalk starts first and calls external function RunGoHttpServer in shared library GoHttpServer. New OS threads and goroutines are created in the Go shared library. These goroutines process HTTP requests. Then Smalltalk calls function GetRequest, which tries to pull request from channel RequestChan. Until some request is actually put to this channel, Smalltalk thread will be blocked. After request received, Go put it to the channel RequestChan. It unblocks Smalltalk thread, function GetRequest writes request's content to the Smalltalk buffer and returns. Then Go waits Smalltalk to process request and push response to channel ResponseChan. Smalltalk does it by calling function PutResponse in library GoHttpServer. When it is done, Go HTTP Server responds to the client.
