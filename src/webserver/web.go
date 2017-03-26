package main

import (
	"net/http"
	"io"
	"log"
)

const listenHost = ""
const listenPort = "4001"

func helloTest(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "hello world\n")
}
func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	errtext := "404 NOT FOUND: " + string(req.Method) + " " + string(req.RequestURI)
	io.WriteString(w, errtext)
	req.ParseForm()
	log.Println(req.Method, req.URL, req.Header, req.RequestURI, req.Form, w.Header())
}

func main() {
	http.HandleFunc("/hello", helloTest)
	http.HandleFunc("/", NotFoundHandler)
	log.Println("Web Server Started, Listen: " + listenPort)
	log.Fatal(http.ListenAndServe(listenHost + ":" + listenPort, nil))
}