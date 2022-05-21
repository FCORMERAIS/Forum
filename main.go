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
				if goodMail(usernameConnect) {
					fmt.Println("entrer2")
					data := connected(usernameConnect, passwordConnect)
				} else {
					fmt.Println("erreur le mail n'est pas bon ")
				}
			} else {
				if passwordGood(strings.Split(password, "")) {
					data := SignUp(email, username, password)
				}
			}
		}
	} else {
		path = "." + path
	}
	http.ServeFile(w, r, path)
}

func Forum(w http.ResponseWriter, r *http.Request) {
	filename := "./templates/Forum.html"
	data := User{
		email:    "flavio@gmail",
		username: "Flavio",
	}
	fmt.Println(data)
	t, _ := template.ParseFiles(filename)
	t.ExecuteTemplate(w, filename, data)
}

func connected(Useremail string, Userpassword string) User {
	db, err := sql.Open("sqlite3", "./BD/Forum.db")
	if err != nil {
		fmt.Println(err)
	}
	resulttest, err2 := db.Prepare("SELECT PasswordHash FROM User WHERE Email = ?")
	testvalue, err3 := resulttest.Query(Useremail)
	if err2 != nil || err3 != nil {
		fmt.Println(err2, testvalue)
	}
	return User{
		email:    Useremail,
		username: Useremail, // A CHANGER !!!!!!!!!!!
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
