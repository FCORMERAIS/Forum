package main

import (
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func SignUp(Useremail string, Userusername string, Userpassword string) User {
	var err error
	Id := uuid.Must(uuid.NewV4(), err)
	fmt.Println(Id) // ID PERMETTANT DE SAVOIR QUI EST-CE A AJOUTER DANS LA ABSE DE DONNÉES APRES
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	result, err := db.Prepare("INSERT INTO User (ID_User,Email, UserName, PasswordHash) VALUES (?,?, ?, ?)") // A CHANGER
	if err != nil {
		fmt.Println(err)
	}
	_, err2 := result.Exec(Id, Useremail, Userusername, Userpassword)
	if err2 != nil {
		fmt.Println(err2)
	}
	db.Close()
	return User{
		Id:       Id.String(),
		Username: Userusername,
	}
}

func goodMail(mail string) string {
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
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
	if err3 != nil {
		db.Close()
		return ""
	}
	for result.Next() {
		result.Scan(&password)
	}
	return password
}

func GetUsernameByID(UUID string) string {
	db, err1 := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err1 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err1)
	}
	statement, err2 := db.Prepare("SELECT UserName FROM User WHERE ID_User = ?")
	if err2 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err2)
	}
	result, err3 := statement.Query(UUID)
	if err3 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err3)
	}
	db.Close()
	var Username string
	for result.Next() {
		result.Scan(&Username)
	}
	return Username
}

func SendPostinDB(message string) {
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println("Erreur ouverture du fichier :")
		fmt.Println(err)
	}
	statement, err := db.Prepare("INSERT INTO Post (ID_Post, ID_User_Post, ID_Catégorie_Post, Text_Post) VALUES (?,?,?,?)")
	var eRR error
	_, err2 := statement.Exec(55, uuid.Must(uuid.NewV4(), eRR), 47, message)
	if err != nil || err2 != nil {
		fmt.Println("Erreur d'insertion :")
		fmt.Println(err)
		fmt.Println(err2)
	}
	db.Close()
}

func GetPostDB() []Post {
	var postList ArrayPosts
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println("Erreur ouverture :")
		fmt.Println(err)
	}
	resulttest, err := db.Query("SELECT ID_User_Post, Text_Post, Like, Dislike FROM Post")
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

func connected(useremail string) User {
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	var ID string
	var username string
	tempo, err2 := db.Prepare("SELECT ID_User,UserName FROM User WHERE Email = ?") // A CHANGER USERNAME AVEC ID
	result, err3 := tempo.Query(useremail)

	if err2 != nil || err3 != nil {
		fmt.Println(err2)
	}
	db.Close()
	for result.Next() {
		result.Scan(&ID, &username)
	}
	return User{
		Id:       ID,
		Username: username,
	}
}