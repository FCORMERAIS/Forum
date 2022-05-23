package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string
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
				cookie, err := r.Cookie("UserSessionId")
				if err != nil {
					cookie = &http.Cookie{
						Name: "UserSessionId",
					}
				} else {
					cookie.MaxAge = -1
				}
				if passwordAccount != "" {
					fmt.Println("entrer2")
					if passwordAccount == passwordConnect {
						data := connected(usernameConnect)
						cookie.Value = data.Username
						cookie.MaxAge = 300
						http.SetCookie(w, cookie)
					} else {
						fmt.Println("LE MDP EST PAS BON ")
					}
				} else {
					fmt.Println("vous n'avez pas rentrer de mot de passe ou le mail n'est pas bon ")
				}
			} else {
				if passwordGood(password, w) {
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
	cookie, err := r.Cookie("UserSessionId")
	if err != nil {
		fmt.Println(err)
	}
	data := User{
		Username: "PAS CONNECTÉ",
	}
	if cookie.Value != "" {
		data.Username = cookie.Value
	}
	fmt.Println(data)
	t, err := template.ParseFiles("./templates/Forum.html", "./templates/header.html")
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	err2 := t.Execute(w, data)
	if err2 != nil {
		fmt.Printf("error2, %s\n", err2)
	}
}

func connected(useremail string) User {
	db, err := sql.Open("sqlite3", "./BD/Forum.db")
	if err != nil {
		fmt.Println(err)
	}
	var username string
	tempo, err2 := db.Prepare("SELECT UserName FROM User WHERE Email = ?")
	result, err3 := tempo.Query(useremail)

	if err2 != nil || err3 != nil {
		fmt.Println(err2)
	}
	db.Close()
	for result.Next() {
		result.Scan(&username)
	}
	return User{
		Username: username, // A CHANGER !!!!!!!!!!!
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
		Username: Userusername,
	}
}

func passwordGood(mdp string, w http.ResponseWriter) bool {
	if len(mdp) < 10 {
		fmt.Fprintf(w, `<p class="error_message">le MDP EST TROP COURT</p>`)
		return false
	}
	r, _ := regexp.Compile("1|2|3|4|5|6|7|8|9")
	if !r.MatchString(mdp) {
		fmt.Fprintf(w, `<p class="error_message">IL MANQUE UN CHIFFRE DANS VOTRE MDP</p>`)
		return false
	}
	r2, _ := regexp.Compile("/|\\|#|*|;|!|?|$|}|[")
	if !r2.MatchString(mdp) {
		fmt.Fprintf(w, `<p class="error_message">IL MANQUE UN CARACTèRE SPÉCIAL</p>`)
		return false
	}
	r3, _ := regexp.Compile("[A-Z]")
	if !r3.MatchString(mdp) {
		fmt.Fprintf(w, `<p class="error_message">il faut une majuscule pour votre mdp</p>`)
		return false
	}
	r4, _ := regexp.Compile("[a-z]")
	if !r4.MatchString(mdp) {
		fmt.Fprintf(w, `<p class="error_message">Erreur il faut une minuscule dans votre mdp</p>`)
		return false
	}
	return true
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
