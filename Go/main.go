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
	http.HandleFunc("/JsonCategories", GetCategories)
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
	var data = User{
		Username: "Invité",
	}
	fmt.Println(r.Method)
	if path == "/" {
		if r.Method == "POST" { // Lorsque que l'on rentre dans le POST cc'est que l'utilisateur veut se connecter ou s'enregistrer
			var usernameRegister = r.FormValue("username")
			var passwordRegister = r.FormValue("password")
			var emailRegister = r.FormValue("email")
			var EmailConnect = r.FormValue("email2")
			var passwordConnect = r.FormValue("password2")
			if EmailConnect != "" && passwordConnect != "" { // l'utilisateur essaye de se connecter
				passwordAccount := goodMail(EmailConnect)
				if passwordAccount != "" {
					if CheckPasswordHash(passwordConnect, passwordAccount) {
						cookie := &http.Cookie{
							Name: "UserSessionId",
						}
						data = connected(EmailConnect)
						cookie.Value = data.Id
						cookie.MaxAge = 300
						http.SetCookie(w, cookie)
					}
				} else {
					fmt.Println("Le mot de passe n'est pas bon ")
				}
			} else if passwordRegister != "" && usernameRegister != "" && emailRegister != "" { // l'utilisateur essaye de s'enregistrer
				if passwordGood(passwordRegister, w) {
					passwordHash, err := HashPassword(passwordRegister)
					if err != nil {
						fmt.Println(err)
					}
					data = SignUp(emailRegister, usernameRegister, passwordHash)
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
			} else if usernameRegister == "" || passwordConnect == "" || emailRegister == "" || passwordRegister == "" || EmailConnect == "" { // l'utilisateur essaye de se déconnecté
				data.Username = "Invité"
			} else { // sinon il y a un roblème on affiche la page ERROR 404
				t, err := template.ParseFiles("../templates/error404.html")
				if err != nil {
					fmt.Println(err)
				}
				err2 := t.Execute(w, data)
				if err2 != nil {
					fmt.Println(err2)
				}
			}
		}
	} else {
		path = ".." + path
	}
	if r.URL.Path == "/" {
		t, err := template.ParseFiles("../templates/server.html", "../templates/header.html")
		if err != nil {
			fmt.Println(err)
		}
		err2 := t.Execute(w, data)
		if err2 != nil {
			fmt.Println(err2)
		}
	} else {
		http.ServeFile(w, r, path)
	}
}

func Forum(w http.ResponseWriter, r *http.Request) {
	var Categories = GetAllCategories()
	var Page ForumPage
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
	if r.Method == "POST" && r.FormValue("Message_Value") != "" && data.Username != "Invité" {
		SendPostinDB(r.FormValue("Message_Value"), data.Id)

		fmt.Println(r.FormValue("Categorie"))
	}
	t, err := template.ParseFiles("../templates/Forum.html")
	if err != nil {
		fmt.Println(err)
	}
	Page.User = data
	Page.ListCategories = Categories
	err2 := t.Execute(w, Page)
	if err2 != nil {
		fmt.Println(err2)
	}
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetAllCategories())
}
