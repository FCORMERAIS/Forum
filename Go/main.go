package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var Port = "127.0.0.1:5555"

func main() {
	fileserver := http.FileServer(http.Dir("static"))
	http.Handle("/static", http.StripPrefix("/static", fileserver))
	http.HandleFunc("/", Acceuil)
	http.HandleFunc("/Forum", Forum)
	http.HandleFunc("/donneesJson", GetJson)
	fmt.Println("Serving @ : ", "http://"+Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}

func GetJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetPostDB())
}
func Acceuil(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	var data User
	fmt.Println(r.Method)
	if path == "/" {
		path = "../templates/server.html"
		if r.Method == "POST" {
			var username = r.FormValue("username")
			var password = r.FormValue("password")
			var email = r.FormValue("email")
			var EmailConnect = r.FormValue("email2")
			var passwordConnect = r.FormValue("password2")
			if EmailConnect != "" && passwordConnect != "" {
				fmt.Println("entrer1")
				passwordAccount := goodMail(EmailConnect)
				cookie, err := r.Cookie("UserSessionId")
				if err != nil {
					cookie = &http.Cookie{
						Name: "UserSessionId",
					}
				} else {
					cookie.MaxAge = -1
					data.Username = "Invité"
				}
				if passwordAccount != "" {
					fmt.Println("entrer2")
					if CheckPasswordHash(passwordConnect, passwordAccount) {
						data := connected(EmailConnect)
						cookie.Value = data.Id
						cookie.MaxAge = 300
						http.SetCookie(w, cookie)
					}
				} else {
					fmt.Println("vous n'avez pas rentrer de mot de passe ou le mail n'est pas bon ")
				}
			} else {
				if passwordGood(password, w) {
					passwordHash, err := HashPassword(password)
					if err != nil {
						fmt.Println(err)
					}
					data := SignUp(email, username, passwordHash)
					cookie := &http.Cookie{
						Name: "UserSessionId",
					}
					cookie.Value = data.Id
					cookie.MaxAge = 300
					http.SetCookie(w, cookie)
				} else {
					data.Username = "Invité"
					fmt.Fprintf(w, "UNE ERREUR EST SURVENUE ")
				}
			}
		}
	} else {
		path = ".." + path
		cookie, err := r.Cookie("UserSessionId")
		if err != nil {
			cookie = &http.Cookie{
				Name: "UserSessionId",
			}
			data.Username = "Invité"
		} else {
			data.Username = GetUsernameByID(cookie.Value)
		}
	}
	if path == "../static/style.css" || path == "../images/Background.png" || path == "../js/index.js" {
		http.ServeFile(w, r, path)
	} else {
		t, err := template.ParseFiles("../templates/server.html")
		if err != nil {
			fmt.Println(err)
		}
		err2 := t.Execute(w, data)
		if err2 != nil {
			fmt.Println(err2)
		}
	}
}

func Forum(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	cookie, err := r.Cookie("UserSessionId")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "UserSessionId",
			Value: "Invité",
		}
	}
	cookie.MaxAge = 300
	data := User{
		Id: cookie.Value,
	}
	if data.Id != "Invité" {
		username := GetUsernameByID(data.Id)
		data.Username = username
	} else {
		data.Username = "Invité"
	}
	t, err := template.ParseFiles("../templates/Forum.html")
	if err != nil {
		fmt.Printf("error %s \n", err)
	}
	if r.Method == "POST" {
		SendPostinDB(r.FormValue("SendPost"))
	}
	err2 := t.Execute(w, data)
	if err2 != nil {
		fmt.Printf("error2, %s\n", err2)
	}
}
