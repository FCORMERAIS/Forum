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

func SendPostinDB(message string, Id_User string, categorie string) {
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println("Erreur ouverture du fichier :")
		fmt.Println(err)
	}
	statement, err := db.Prepare("INSERT INTO Post (ID_Post, ID_User_Post, ID_Catégorie_Post, Text_Post) VALUES (?,?,?,?)")
	var eRR error
	_, err2 := statement.Exec(uuid.Must(uuid.NewV4(), eRR), Id_User, categorie, message)
	if err != nil || err2 != nil {
		fmt.Println("Erreur d'insertion :")
		fmt.Println(err)
		fmt.Println(err2)
	}
	db.Close()
}

func GetPostDB(filter string) []Post {
	var postList []Post
	var resultPost *sql.Rows
	var ID_Categorie_Filtre string
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println("Erreur ouverture :")
		fmt.Println(err)
	}
	if filter == "" {
		resultPost, err = db.Query("SELECT ID_Post, ID_User_Post, ID_Catégorie_Post, Text_Post, Like, Dislike FROM Post")
	} else {
		prepareforRecupID_Categorie, err2 := db.Prepare("SELECT ID_Categorie FROM Categorie WHERE Name = ?")
		resulte_ID, err := prepareforRecupID_Categorie.Query(filter)
		if err != nil || err2 != nil {
			fmt.Println("Erreur de recherche :")
			fmt.Println(err)
			fmt.Println(err2)
		}
		for resulte_ID.Next() {
			resulte_ID.Scan(&ID_Categorie_Filtre)
			prepare, _ := db.Prepare("SELECT ID_Post, ID_User_Post, ID_Catégorie_Post, Text_Post, Like, Dislike FROM Post WHERE ID_Catégorie_Post = ?")
			resultPost, err = prepare.Query(ID_Categorie_Filtre)
		}

	}

	if err != nil {
		fmt.Println("Erreur de recherche :")
		fmt.Println(err)
	}
	var Username string
	var Text_Post string
	var id_post string
	var Like string
	var Dislike string
	var numberLike int
	var numberDislike int
	var IdUser string
	var ID_Categorie string
	var CategorieColor string
	var CategorieName string
	var singlePost Post
	for resultPost.Next() {
		Like = ""
		Dislike = ""
		resultPost.Scan(&id_post, &IdUser, &ID_Categorie, &Text_Post, &Like, &Dislike)
		fmt.Println(Like, id_post)
		Username = GetUsernameByID(IdUser)
		numberLike = KnowLike(Like)
		numberDislike = KnowLike(Dislike)
		resultCategorie, err := db.Prepare("SELECT Name, Color FROM Categorie WHERE ID_Categorie = ?")
		result, err := resultCategorie.Query(ID_Categorie)
		if err != nil {
			fmt.Println("Erreur de recherche :")
			fmt.Println(err)
		}
		for result.Next() {
			result.Scan(&CategorieName, &CategorieColor)
		}
		singlePost = Post{
			Username:       Username,
			TextPost:       Text_Post,
			LikePost:       numberLike,
			DislikePost:    numberDislike,
			IdPost:         id_post,
			CategorieColor: CategorieColor,
			CategorieName:  CategorieName,
		}
		postList = append(postList, singlePost)
	}
	db.Close()
	return postList
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

func PostWithCategories(categorie string) []Post {
	var Posts []Post
	return Posts
}

func GetAllCategories() []Categorie {
	var categories []Categorie
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	tableCategorie, err2 := db.Query("SELECT Name, Color FROM Categorie")
	if err2 != nil {
		fmt.Println(err2)
	}
	var url string
	var name string
	var color string
	for tableCategorie.Next() {
		tableCategorie.Scan(&name, &color)
		url = "/Forum#" + name
		categories = append(categories, Categorie{URL: url, Name: name, Color: color})
	}
	return categories
}

func addUserLike(userID string, post_ID string) {
	var resultPost string
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, err2 := db.Prepare("SELECT Like FROM Post WHERE ID_Post = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	result, err3 := statement.Query(post_ID)
	if err3 != nil {
		fmt.Println(err3)
	}
	for result.Next() {
		result.Scan(&resultPost)
	}
	db.Close()
	resultPost += userID + " "
	db2, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement2, err2 := db2.Prepare("UPDATE Post SET Like = ? WHERE ID_Post = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	tempo, err3 := statement2.Exec(resultPost, post_ID)
	if err3 != nil {
		fmt.Println(err3)
	}
	if tempo == nil {
		fmt.Println("tempo is empty")
	}
	db2.Close()
}

func addUserDislike(userID string, post_ID string) {
	var resultPost string
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, err2 := db.Prepare("SELECT Dislike FROM Post WHERE ID_Post = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	result, err3 := statement.Query(post_ID)
	if err3 != nil {
		fmt.Println(err3)
	}
	for result.Next() {
		result.Scan(&resultPost)
	}
	db.Close()
	resultPost += userID + " "
	db2, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement2, err2 := db2.Prepare("UPDATE Post SET Dislike = ? WHERE ID_Post = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	tempo, err3 := statement2.Exec(resultPost, post_ID)
	if err3 != nil {
		fmt.Println(err3)
	}
	if tempo == nil {
		fmt.Println("tempo is empty")
	}
	db2.Close()
}

func GetPostLike(uuid string) string {
	var likestr string
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	tableCategorie, err2 := db.Prepare("SELECT Like FROM Post WHERE ID_Post = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	result, err3 := tableCategorie.Query(uuid)
	if err3 != nil {
		fmt.Println(err3)
	}
	for result.Next() {
		result.Scan(&likestr)
	}
	db.Close()
	return likestr
}

func GetPostDisike(uuid string) string {
	var dislikestr string
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	tableCategorie, err2 := db.Prepare("SELECT Dislike FROM Post WHERE ID_Post = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	result, err3 := tableCategorie.Query(uuid)
	if err3 != nil {
		fmt.Println(err3)
	}
	for result.Next() {
		result.Scan(&dislikestr)
	}
	db.Close()
	return dislikestr
}

func GetIdCategorie(categorie string) string {
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, err2 := db.Prepare("SELECT ID_Categorie FROM Categorie WHERE Name = ?")
	if err2 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err2)
	}
	result, err3 := statement.Query(categorie)
	if err3 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err3)
	}
	db.Close()
	var IdCategorie string
	for result.Next() {
		result.Scan(&IdCategorie)
	}
	return IdCategorie

}
