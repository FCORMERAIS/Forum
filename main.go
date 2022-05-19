package main

import (
	"fmt"
	"log"
	"net/http"
)

var Port = "127.0.0.1:5555"

func main() {
	http.HandleFunc("/", ServeFiles)
	fmt.Println("Serving @ : ", "http://"+Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}

func ServeFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	if path == "/" {
		path = "./templates/server.html"
	} else if path == "/Forum" {
		path = "./templates/Forum.html"
	} else {
		path = "." + path
	}
	http.ServeFile(w, r, path)
}
