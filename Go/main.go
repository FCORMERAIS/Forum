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
var UserPost []Post
var UserLikePost []Post

func main() {
	fileserver := http.FileServer(http.Dir("static"))
	http.Handle("/static", http.StripPrefix("/static", fileserver))
	http.HandleFunc("/", testPath)                    //si on tombe sur une URL non attribué on tests le chemin d'accées (pour verifier si il s'agit de js ou de css) sinon on affiche un message d'erreur
	http.HandleFunc("/Acceuil", Acceuil)              // l'url /Acceuil renvoie vers la page de connection (il s'agit de la page sur lequel l'utilisateur est censé arriver en premier)
	http.HandleFunc("/Forum", Forum)                  //cette url ramène versle forum ou il y a les posts
	http.HandleFunc("/Post", GetPost)                 // GetJson permet de stocker les posts utile dans un ficher JSON en fonction des demandes de l'utilisateur
	http.HandleFunc("/JsonCategories", GetCategories) //GetCategories fonctionne comme Json mais on y stock des Categories
	fmt.Println("Serving @ : ", "http://"+Port+"/Acceuil")
	log.Fatal(http.ListenAndServe(Port, nil))
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetAllCategories())
}

//fonctoin permettant de recuperer tout les posts, il renvoie une liste de posts avec l'objet Post
func GetPost(w http.ResponseWriter, r *http.Request) {
	var data User
	cookie, err := r.Cookie("UserSessionId")
	if err != nil {
		data.Username = "Invité"
		data.Id = ""
	} else {
		data.Id = cookie.Value
	}
	w.Header().Set("Content-Type", "application/json")
	if UserPost != nil {
		json.NewEncoder(w).Encode(UserPost)
	} else if UserLikePost != nil {
		json.NewEncoder(w).Encode(UserLikePost)
	} else {
		json.NewEncoder(w).Encode(GetPostDB(filter, data.Id))
	}
}

//ici on test les URL pour voir si elle peuvent etre utile si elle ne le sont pas on affiche une erreur 404
func testPath(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/js/Forum.js" || r.URL.Path == "/js/index.js" || r.URL.Path == "/static/style_main-page.css" || r.URL.Path == "/static/style.css" || r.URL.Path == "/images/Background.png" {
		http.ServeFile(w, r, ".."+r.URL.Path)
	} else {
		error404(w, r)
	}
}

//Acceuil est la page d'acceuil de notre site cest ici que l'utilisateur pourra se creer un compte ou bien se connecter ou se deconnecter
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
						cookie.MaxAge = 3600 * 24
						http.SetCookie(w, cookie)
					}
				} else {
					fmt.Fprintf(w, "Le mot de passe n'est pas bon ")
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
					cookie.MaxAge = 3600 * 24
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

//fonction forum l'utilisateur pourra liker des posts en poster ou bien simplement les regarder si il n'est pas connecter
func Forum(w http.ResponseWriter, r *http.Request) {
	var ERROR bool = false
	fmt.Println(r.URL.Path)
	cookie, err := r.Cookie("UserSessionId")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "UserSessionId",
			Value: "Invité",
		}
	}
	cookie.MaxAge = 3600 * 24
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
		if r.FormValue("Message_Value") != "" && data.Username != "Invité" { // ici l'utilisateur essaye de poster un message
			SendPostinDB(r.FormValue("Message_Value"), data.Id, GetIdCategorie(r.FormValue("Categorie")))
		} else if r.FormValue("Dislike") != "" && data.Username != "Invité" { // l'utilsateur essaye de mettre un dislike
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
		} else if data.Username != "Invité" && r.FormValue("LikeComm") != "" { // on test si l'utilisateur essaye de liker un messge et si il est connecté
			if !Like(GetCommentLike(r.FormValue("LikeComm")), data.Id) { // l'utilisateur like un commentaire qu'il n'a pas encore liker
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
		} else if r.FormValue("SeeAllPost") != "" { // si l'utilisateur appuie sur le boutton voir tout les posts on renitialise tous les Filtres
			UserPost = nil
			UserLikePost = nil
			filter = ""
		} else if r.FormValue("categorieForm") != "" { // on change la valeur de de filter pour qu'il soit associé a la
			UserPost = nil
			UserLikePost = nil
			filter = r.FormValue("categorieForm")
		} else if r.FormValue("SeeOurPost") != "" { // ici l'utilisateur veut voir ses posts donc on va les cherchers avec GetUserPost
			UserLikePost = nil
			filter = ""
			UserPost = GetUserPost(r.FormValue("SeeOurPost"))
		} else if r.FormValue("SeePostLike") != "" { //  ici l'utilisateur veut voir ses posts liker donc on va les cherchers avec GetUserPostLike
			UserPost = nil
			filter = ""
			UserLikePost = GetUserPostLike(r.FormValue("SeePostLike"))
		} else if r.FormValue("delete") != "" { // ici l'utilisaeur souhaite supprimer un de ses posts
			DeletePost(r.FormValue("delete"))
		} else if r.FormValue("nameCategorieAdd") != "" && r.FormValue("colorAddCategorie") != "" {
			fmt.Println(r.FormValue("nameCategorieAdd"))
			fmt.Println(r.FormValue("colorAddCategorie"))
			AddCategorie(r.FormValue("nameCategorieAdd"), r.FormValue("colorAddCategorie"))
		} else { // sinon il y a une erreur et lance l'erreur 500
			error500(w, r)
			ERROR = true
		}
	}
	if !ERROR {
		t, err := template.ParseFiles("../templates/Forum.html") // on charge la templates du Forum
		if err != nil {
			fmt.Println(err)
		}
		err2 := t.Execute(w, data) // on l'éxecute
		if err2 != nil {
			fmt.Println(err2)
			error500(w, r)
		}
	}
}

//fonction permettant d'afficher a l'utilisateur qu'il s'est tromper d'URL et que c'est une erreur 404
func error404(w http.ResponseWriter, r *http.Request) { // fonction qui affiche la page de l'erreur 404
	tmpl, err := template.ParseFiles("../templates/error404.html") // utilisation du fichier error pour le template
	if err != nil {
		fmt.Println(err)
	}
	tmpl.Execute(w, nil) // exécute le template sur la page html
}

//fonction qui permet d'afficher les erreur 500 pour que l'utilisateur comprenne qu'il y a une erreur interne
func error500(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../templates/error500.html") // utilisation du fichier error pour le template
	if err != nil {
		fmt.Println(err)
	}
	tmpl.Execute(w, nil) // exécute le template sur la page html
}
