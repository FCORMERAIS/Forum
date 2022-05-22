package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	email    string
	username string
}

var Port = ":5555"

func main() {
	fileserver := http.FileServer(http.Dir("static"))
	http.Handle("/static", http.StripPrefix("/static", fileserver))
	http.HandleFunc("/", Acceuil)
	http.HandleFunc("/Forum", Forum)
	fmt.Println("Serving @ : ", "http://127.0.0.1"+Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}

func Acceuil(w http.ResponseWriter, r *http.Request) {
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
				passwordAccount := goodMail(usernameConnect)
				if passwordAccount != "" {
					fmt.Println("entrer2")
					if passwordAccount == passwordConnect {
						fmt.Println("CEST GOOOOOD")
						data := connected(usernameConnect)
						fmt.Println(data)
					} else {
						fmt.Println("LE MDP EST PAS BON ")
					}
				} else {
					fmt.Println("vous n'avez pas rentrer de mot de passe ou le mail n'est pas bon ")
				}
			} else {
				if passwordGood(strings.Split(password, "")) {
					data := SignUp(email, username, password)
					fmt.Println(data)
				}
			}
		}
	} else {
		path = "." + path
	}
	http.ServeFile(w, r, path)
}

func Forum(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	data := User{
		email:    "flavio@gmail",
		username: "Flavio",
	}
	fmt.Println(data)
	t, _ := template.ParseFiles("./templates/Forum.html")
	t.ExecuteTemplate(w, "./templates/Forum.html", data)
}

func connected(Useremail string) User {
	db, err := sql.Open("sqlite3", "./BD/Forum.db")
	if err != nil {
		fmt.Println(err)
	}
	var username string
	tempo, err2 := db.Prepare("SELECT UserName FROM User WHERE Email = ?")
	result, err3 := tempo.Query(Useremail)

	if err2 != nil || err3 != nil {
		fmt.Println(err2)
	}
	db.Close()
	for result.Next() {
		result.Scan(&username)
	}
	return User{
		email:    Useremail,
		username: username, // A CHANGER !!!!!!!!!!!
	}
}

func SignUp(Useremail string, Userusername string, Userpassword string) User {
	db, err := sql.Open("sqlite3", "./BD/Forum.db")
	if err != nil {
		fmt.Println(err)
	}
	result, err := db.Prepare("INSERT INTO User (Email, UserName, PasswordHash) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println(err)
	}
	_, err2 := result.Exec(Useremail, Userusername, Userpassword)
	if err2 != nil {
		fmt.Println(err2)
	}
	db.Close()
	return User{
		email:    Useremail,
		username: Userusername,
	}
}

func passwordGood(mdp []string) bool {
	return len(mdp) > 10
}

func goodMail(mail string) string {
	db, err := sql.Open("sqlite3", "./BD/Forum.db")
	if err != nil {
		db.Close()
		return ""
	}
	resulttest, err2 := db.Prepare("SELECT PasswordHash FROM User WHERE Email = ?")
	if err2 != nil {
		db.Close()
		return ""
	}
	var password string
	result, err3 := resulttest.Query(mail)
	db.Close()
	for result.Next() {
		result.Scan(&password)
	}
	if err3 != nil {
		db.Close()
		return ""
	}
	return password
}
