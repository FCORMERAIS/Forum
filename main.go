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
		// tmpl, err := template.ParseFiles("./templates/Forum.html", "./templates/header.html", "./templates/server.html")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		fmt.Println(r.Method)
		if r.Method == "POST" {
			fmt.Println(r.FormValue("registerForm"))
		}
		// tmpl.ExecuteTemplate(w, "listartists", nil)

	} else if path == "/Forum" {
		path = "./templates/Forum.html"
	} else {
		path = "." + path
	}
	http.ServeFile(w, r, path)
}
