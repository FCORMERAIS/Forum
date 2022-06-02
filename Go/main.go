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
var filter string = ""

func main() {
	fileserver := http.FileServer(http.Dir("static"))
	http.Handle("/static", http.StripPrefix("/static", fileserver))
	http.HandleFunc("/", testPath)
	http.HandleFunc("/Acceuil", Acceuil)
	http.HandleFunc("/Forum", Forum)
	http.HandleFunc("/Post", GetJson)
	http.HandleFunc("/JsonCategories", GetCategories)
	fmt.Println("Serving @ : ", "http://"+Port+"/Acceuil")
	log.Fatal(http.ListenAndServe(Port, nil))
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetAllCategories())
}

func GetJson(w http.ResponseWriter, r *http.Request) {
	var data User
	cookie, err := r.Cookie("UserSessionId")
	if err != nil {
		data.Username = "Invité"
		data.Id = ""
	} else {
		data.Id = cookie.Value
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetPostDB(filter, data.Id))
}

func testPath(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/js/Forum.js" || r.URL.Path == "/js/index.js" || r.URL.Path == "/static/style_main-page.css" || r.URL.Path == "/static/style.css" || r.URL.Path == "/images/Background.png" {
		http.ServeFile(w, r, ".."+r.URL.Path)
	} else {
		error404(w, r)
	}
}

func Acceuil(w http.ResponseWriter, r *http.Request) {
	filter = ""
	path := r.URL.Path
	fmt.Println(path)
	var data User
	cookie, err := r.Cookie("UserSessionId")
	if err != nil {
		data.Username = "Invité"
	} else {
		data.Username = GetUsernameByID(cookie.Value)
	}
	fmt.Println(r.Method)
	if path == "/Acceuil" {
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
						error500(w, r)
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
				error500(w, r)
			}
		}
	}
	if r.URL.Path == "/Acceuil" {
		t, err := template.ParseFiles("../templates/server.html", "../templates/header.html") // on parse les fichier html que l'on a besoin pour afficher la page voulut
		if err != nil {
			fmt.Println(err)
			error500(w, r) // si il y a un problème on lance la fonctio nerror 404
		}
		err2 := t.Execute(w, data)
		if err2 != nil {
			fmt.Println(err2)
			error500(w, r)
		}
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
	if r.Method == "POST" {
		if r.FormValue("Message_Value") != "" && data.Username != "Invité" {
			SendPostinDB(r.FormValue("Message_Value"), data.Id, GetIdCategorie(r.FormValue("Categorie")))
		} else if r.FormValue("Dislike") != "" && data.Username != "Invité" {
			if !Like(GetPostDisike(r.FormValue("Dislike")), data.Id) {
				deleteUserLikePost(data.Id, r.FormValue("Dislike"))
				addUserDislikePost(data.Id, r.FormValue("Dislike"))
			} else { // ici l'utilisateur essaie de disliker un post mais il la déjà disliker
				deleteUserLikePost(data.Id, r.FormValue("Dislike"))
				deleteUserDislikePost(data.Id, r.FormValue("Dislike"))
			}
		} else if r.FormValue("Like") != "" && data.Username != "Invité" { // ici lutilisateur essaye de liker un post quil na pas encore liker
			if !Like(GetPostLike(r.FormValue("Like")), data.Id) {
				deleteUserDislikePost(data.Id, r.FormValue("Like"))
				addUserLikePost(data.Id, r.FormValue("Like"))
			} else { // ici l'utilisateur essaye de liker un post qu'il a déjà liker
				deleteUserDislikePost(data.Id, r.FormValue("Like"))
				deleteUserLikePost(data.Id, r.FormValue("Like"))
			}
		} else if r.FormValue("idPost") != "" && data.Username != "Invité" && r.FormValue("textCommentary") != "" { // l'utilisateur post un commentaire
			addCommentary(r.FormValue("idPost"), r.FormValue("textCommentary"), data.Id)
		} else if data.Username != "Invité" && r.FormValue("LikeComm") != "" { // l'utilisateur like un commentaire qu'il n'a pas encore liker
			if !Like(GetCommentLike(r.FormValue("LikeComm")), data.Id) {
				deleteUserDislikeComment(data.Id, r.FormValue("LikeComm"))
				addUserLikeComment(data.Id, r.FormValue("LikeComm"))
			} else { // l'utilisateur like un commentaire qu'il a déjà liker
				deleteUserDislikeComment(data.Id, r.FormValue("LikeComm"))
				deleteUserLikeComment(data.Id, r.FormValue("LikeComm"))
			}
		} else if data.Username != "Invité" && r.FormValue("DislikeComm") != "" { // l'utilisateur dislike un commentaire qu'il n'a pas encore liker
			if !Like(GetCommentDislike(r.FormValue("DislikeComm")), data.Id) {
				deleteUserLikeComment(data.Id, r.FormValue("DislikeComm"))
				addUserDislikeComment(data.Id, r.FormValue("DislikeComm"))
			} else { // l'utilisateur dislike un commentaire u'il a déjà disliker
				deleteUserDislikeComment(data.Id, r.FormValue("DislikeComm"))
				deleteUserLikeComment(data.Id, r.FormValue("DislikeComm"))
			}
		} else if r.FormValue("categorieForm") != "" {
			filter = r.FormValue("categorieForm")
		} else { // sinon il y a une erreur et lance l'erreur 404
			error500(w, r)
		}
	}
	t, err := template.ParseFiles("../templates/Forum.html") // on charge la templates du Forum
	if err != nil {
		fmt.Println(err)
	}
	Page.User = data
	Page.ListCategories = Categories
	err2 := t.Execute(w, Page) // on l'éxecute
	if err2 != nil {
		fmt.Println(err2)
		error500(w, r)
	}
}

func error404(w http.ResponseWriter, r *http.Request) { // fonction qui affiche la page de l'erreur 404
	tmpl, err := template.ParseFiles("../templates/error404.html") // utilisation du fichier error pour le template
	if err != nil {
		fmt.Println(err)
	}
	tmpl.Execute(w, nil) // exécute le template sur la page html
}

func error500(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../templates/error500.html") // utilisation du fichier error pour le template
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, "hello")
	tmpl.Execute(w, nil) // exécute le template sur la page html
}
