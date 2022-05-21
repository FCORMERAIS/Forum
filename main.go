package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var Port = ":5555"

func main() {
	http.HandleFunc("/", ServeFiles)
	fmt.Println("Serving @ : ", "http://127.0.0.1"+Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}

func ServeFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	fmt.Println(r.Method)
	if path == "/" {
		path = "./templates/server.html"
		if r.Method == "POST" {
			var username = r.FormValue("username")
			var password = r.FormValue("password")
			var email = r.FormValue("email")
			var usernameConnect = r.FormValue("email2")
			var passwordConnect = r.FormValue("password2")
			if usernameConnect != "" && passwordConnect != "" {
				fmt.Println("entrer1")
				if goodMail(usernameConnect) {
					fmt.Println("entrer2")
					db, err := sql.Open("sqlite3", "./BD/Forum.db")
					if err != nil {
						fmt.Println(err)
					}
					resulttest, err2 := db.Prepare("SELECT PasswordHash FROM User WHERE Email = ?")
					testvalue, err3 := resulttest.Query(usernameConnect)
					if err2 != nil || err3 != nil {
						fmt.Println(err2, testvalue)
					}
					// fmt.Println(usernameConnect, testvalue)
				} else {
					fmt.Println("erreur le mail n'est pas bon ")
				}
			} else {
				db, err := sql.Open("sqlite3", "./BD/Forum.db")
				if err != nil {
					fmt.Println(err)
				}
				result, err := db.Prepare("INSERT INTO User (Email, UserName, PasswordHash) VALUES (?, ?, ?)")
				if err != nil {
					fmt.Println(err)
				}
				_, err2 := result.Exec(email, username, password)
				if err2 != nil {
					fmt.Println(err2)
				}
				db.Close()
				if passwordGood(strings.Split(password, "")) {
					path = "./templates/Forum.html"
				}
			}
		}
	} else if path == "/Forum" {
		path = "./templates/Forum.html"
	} else {
		path = "." + path
	}
	http.ServeFile(w, r, path)
}

func Forum(w http.ResponseWriter, r *http.Request) {

}

func passwordGood(mdp []string) bool {
	if len(mdp) < 10 {
		return false
	}
	return true
}

func goodMail(mail string) bool {
	db, err := sql.Open("sqlite3", "./BD/Forum.db")
	if err != nil {
		db.Close()
		return false
	}
	resulttest, err2 := db.Prepare("SELECT UserName FROM User WHERE Email = ?")
	if err2 != nil {
		db.Close()
		return false
	}
	result, err3 := resulttest.Query(mail)
	fmt.Println(resulttest, result)
	db.Close()
	if err3 != nil {
		db.Close()
		return false
	}
	return false
}
