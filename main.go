package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
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

		// tmpl.ExecuteTemplate(w, "listartists", nil)

	} else if path == "/Forum" {
		path = "./templates/Forum.html"
		if r.Method == "POST" {
			SendPostinDB(r.FormValue("SendPost"))

		}
	} else {
		path = "." + path
	}
	http.ServeFile(w, r, path)
}

func SendPostinDB(message string) {
	db, err := sql.Open("sqlite3", "./BD/Forum.db")
	if err != nil {
		fmt.Println("Erreur ouverture du fichier :")
		fmt.Println(err)
	}
	statement, err := db.Prepare("INSERT INTO Post (ID_Post, ID_User, ID_Catégorie, Text_Post) VALUES (?,?,?,?)")
	_, err2 := statement.Exec(2206, rand.Int(), 2306, message)
	if err != nil || err2 != nil {
		fmt.Println("Erreur d'insertion :")
		fmt.Println(err)
	}
	db.Close()
}
