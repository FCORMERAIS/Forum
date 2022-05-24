package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

//--------------------------------------------------------------------------------
type Post struct {
	IDUser      int
	TextPost    string
	LikePost    bool
	DislikePost bool
}

type ArrayPosts struct {
	arrayPosts []Post
}

//--------------------------------------------------------------------------------

var Port = "127.0.0.1:5555"

func main() {
	http.HandleFunc("/", ServeFiles)
	http.HandleFunc("/donneesJson", GetJson)
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
	_, err2 := statement.Exec(rand.Int(), rand.Int(), rand.Int(), message)
	if err != nil || err2 != nil {
		fmt.Println("Erreur d'insertion :")
		fmt.Println(err)
	}
	db.Close()
}

func GetJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetPostDB())
}

func GetPostDB() []Post {
	var postList ArrayPosts
	db, err := sql.Open("sqlite3", "./BD/Forum.db")
	if err != nil {
		fmt.Println("Erreur ouverture :")
		fmt.Println(err)
	}
	resulttest, err := db.Query("SELECT ID_User, Text_Post, Like, Dislike FROM Post")
	if err != nil {
		fmt.Println("Erreur de recherche :")
		fmt.Println(err)
	}
	var ID_User int
	var Text_Post string
	var Like bool
	var Dislike bool
	var singlePost Post
	for resulttest.Next() {
		resulttest.Scan(&ID_User, &Text_Post, &Like, &Dislike)
		singlePost = Post{
			IDUser:      ID_User,
			TextPost:    Text_Post,
			LikePost:    Like,
			DislikePost: Dislike,
		}
		postList.arrayPosts = append(postList.arrayPosts, singlePost)
	}
	resulttest.Close()
	return postList.arrayPosts
}
