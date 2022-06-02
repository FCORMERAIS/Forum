package main

import (
	"database/sql"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// fonction permettant de rajouter un compte dans la base de données
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

//fonction permettant de vérifier si l'adresse email rentrer par l'utilisateur lorsque qu'il se connecte existe bien si elle existe on renvoie son mot de passe
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

//fonction permettant de recuperer un username depuis un id de User
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

func GetColorCategoryById(UUID string) string {
	db, err1 := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err1 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err1)
	}
	statement, err2 := db.Prepare("SELECT Color FROM Categorie WHERE ID_Categorie = ?")
	if err2 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err2)
	}
	result, err3 := statement.Query(UUID)
	if err3 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err3)
	}
	db.Close()
	var color string
	for result.Next() {
		result.Scan(&color)
	}
	return color
}

func GetNameCategoryById(UUID string) string {
	db, err1 := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err1 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err1)
	}
	statement, err2 := db.Prepare("SELECT Name FROM Categorie WHERE ID_Categorie = ?")
	if err2 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err2)
	}
	result, err3 := statement.Query(UUID)
	if err3 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err3)
	}
	db.Close()
	var NameCategory string
	for result.Next() {
		result.Scan(&NameCategory)
	}
	return NameCategory
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

func GetPostDB(filter string, UserId string) []Post {
	var postList []Post
	var resultPost *sql.Rows
	var ID_Categorie_Filtre string
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println("Erreur ouverture :", err)
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
	db.Close()
	var Like, ID_Categorie, Dislike string
	var singlePost Post
	for resultPost.Next() {
		Like = ""
		Dislike = ""
		resultPost.Scan(&singlePost.IdPost, &singlePost.Username, &ID_Categorie, &singlePost.TextPost, &Like, &Dislike)
		if singlePost.Username == UserId {
			singlePost.SamePersonWhithSession = true
		} else {
			singlePost.SamePersonWhithSession = false
		}
		singlePost.Username = GetUsernameByID(singlePost.Username)
		singlePost.LikePost = KnowLike(Like)
		singlePost.DislikePost = KnowLike(Dislike)
		singlePost.CategorieColor = GetColorCategoryById(ID_Categorie)
		singlePost.CategorieName = GetNameCategoryById(ID_Categorie)
		singlePost.CommentaryPost = GetCommmentary(singlePost.IdPost)
		postList = append(postList, singlePost)
	}
	return postList
}

//fonction permettant de récuperer un Objet User grace au username de la personne
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

//fonction permettant de recuperer toutes les catégories stockés dans la base de données
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

//fonction permettant de rajouter un like sur un commentaire grace au post_id et a l'id de l'utilisateur
func addUserLikePost(userID string, post_ID string) {
	resultPost := GetPostLike(post_ID) + userID + "#"
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

//fonction permettant de rajouter un like sur un commentaire grace au post_id et a l'id de l'utilisateur
func addUserLikeComment(userID string, post_ID string) {
	resultPost := GetCommentLike(post_ID) + userID + "#"
	db2, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement2, err2 := db2.Prepare("UPDATE Commentaire SET Like = ? WHERE ID_Commentaire = ?")
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

//fonction permettant de récuperer les données des likes d'un post grace a l'id du post dans la base SQL
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

//fonction permettant de supprimer un dislikes sur un Post grace au post_id et a l'id de l'utilisateur
func deleteUserDislikePost(userID string, post_ID string) {
	resultPost := strings.Split(GetPostDisike(post_ID), "#")
	index := -1
	for i := 0; i < len(resultPost); i++ {
		if userID == resultPost[i] {
			index = i
		}
	}
	if index != -1 {
		resultPost[index] = resultPost[len(resultPost)-1]
		resultPost[len(resultPost)-1] = ""
		resultPost = resultPost[:len(resultPost)-1]
	}
	resultPost2 := strings.Join(resultPost, "#")
	db2, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement2, err2 := db2.Prepare("UPDATE Post SET Dislike = ? WHERE ID_Post = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	tempo, err3 := statement2.Exec(resultPost2, post_ID)
	if err3 != nil {
		fmt.Println(err3)
	}
	if tempo == nil {
		fmt.Println("tempo is empty")
	}
	db2.Close()
}

//fonction permettant de supprimer un dislikes sur un commentaire grace au post_id et a l'id de l'utilisateur
func deleteUserDislikeComment(userID string, post_ID string) {
	resultPost := strings.Split(GetCommentDislike(post_ID), "#")
	index := -1
	for i := 0; i < len(resultPost); i++ {
		if userID == resultPost[i] {
			index = i
		}
	}
	if index != -1 {
		resultPost[index] = resultPost[len(resultPost)-1]
		resultPost[len(resultPost)-1] = ""
		resultPost = resultPost[:len(resultPost)-1]
	}
	resultPost2 := strings.Join(resultPost, "#")
	db2, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement2, err2 := db2.Prepare("UPDATE Commentaire SET Dislike = ? WHERE ID_Commentaire = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	tempo, err3 := statement2.Exec(resultPost2, post_ID)
	if err3 != nil {
		fmt.Println(err3)
	}
	if tempo == nil {
		fmt.Println("tempo is empty")
	}
	db2.Close()
}

//fonction permettant de rajouter un like sur un Post grace au post_id et a l'id de l'utilisateur
func deleteUserLikePost(userID string, post_ID string) {
	resultPost := strings.Split(GetPostLike(post_ID), "#")
	index := -1
	for i := 0; i < len(resultPost); i++ {
		if userID == resultPost[i] {
			index = i
		}
	}
	if index != -1 {
		resultPost[index] = resultPost[len(resultPost)-1]
		resultPost[len(resultPost)-1] = ""
		resultPost = resultPost[:len(resultPost)-1]
	}
	resultPost2 := strings.Join(resultPost, "#")
	db2, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement2, err2 := db2.Prepare("UPDATE Post SET Like = ? WHERE ID_Post = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	tempo, err3 := statement2.Exec(resultPost2, post_ID)
	if err3 != nil {
		fmt.Println(err3)
	}
	if tempo == nil {
		fmt.Println("tempo is empty")
	}
	db2.Close()
}

//fonction permettant de rajouter un like sur un Commentaire grace au post_id et a l'id de l'utilisateur
func deleteUserLikeComment(userID string, post_ID string) {
	resultPost := strings.Split(GetCommentLike(post_ID), "#")
	index := -1
	for i := 0; i < len(resultPost); i++ {
		if userID == resultPost[i] {
			index = i
		}
	}
	if index != -1 {
		resultPost[index] = resultPost[len(resultPost)-1]
		resultPost[len(resultPost)-1] = ""
		resultPost = resultPost[:len(resultPost)-1]
	}
	resultPost2 := strings.Join(resultPost, "#")
	db2, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement2, err2 := db2.Prepare("UPDATE Commentaire SET Like = ? WHERE ID_Commentaire = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	tempo, err3 := statement2.Exec(resultPost2, post_ID)
	if err3 != nil {
		fmt.Println(err3)
	}
	if tempo == nil {
		fmt.Println("tempo is empty")
	}
	db2.Close()
}

//fonction permettant de rajouter un dislikes sur un Post grace au post_id et a l'id de l'utilisateur
func addUserDislikePost(userID string, post_ID string) {
	resultPost := GetPostDisike(post_ID) + userID + "#"
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

//fonction permettant de rajouter un dislikes sur un commentaire grace au post_id et a l'id de l'utilisateur
func addUserDislikeComment(userID string, post_ID string) {
	resultPost := GetCommentDislike(post_ID) + userID + "#"
	db2, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement2, err2 := db2.Prepare("UPDATE Commentaire SET Dislike = ? WHERE ID_Commentaire = ?")
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

//fonction permettant de recupérer les Likes d'un commentaire grace a son id
func GetCommentLike(uuid string) string {
	var likestr string
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	tableCategorie, err2 := db.Prepare("SELECT Like FROM Commentaire WHERE ID_Commentaire = ?")
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

//fonction permettant de recupérer les dislikes d'un commentaire grace a son id
func GetCommentDislike(uuid string) string {
	var likestr string
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	tableCategorie, err2 := db.Prepare("SELECT Dislike FROM Commentaire WHERE ID_Commentaire = ?")
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

//fonction permettant d'avoir tout les dislike d'un post grace a son ID
func GetPostDisike(uuidPost string) string {
	var dislikestr string
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	tableCategorie, err2 := db.Prepare("SELECT Dislike FROM Post WHERE ID_Post = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	result, err3 := tableCategorie.Query(uuidPost)
	if err3 != nil {
		fmt.Println(err3)
	}
	for result.Next() {
		result.Scan(&dislikestr)
	}
	db.Close()
	fmt.Println("ICI", dislikestr)
	return dislikestr
}

//fonction permettant d'ajouter a la base de donnée un commentaire il suffit d'envoyer l'id du post le text que l'on veut ajouter ainsi que l'id du User
func addCommentary(IdPost string, text string, IDUser string) {
	var err error
	IdCommentary := uuid.Must(uuid.NewV4(), err)
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, err2 := db.Prepare("INSERT INTO Commentaire (ID_Commentaire, Text_Commentaire, Post_ID,Username) VALUES (?,?, ?,?)")
	if err2 != nil {
		fmt.Println(err2)
	}
	statement.Exec(IdCommentary, text, IdPost, IDUser)
	db.Close()
}

//cette fonction permet d'avoir la liste des commentaire d'un post en envoyant l'id d'un post
func GetCommmentary(idPost string) []Commentary {
	var ListCommentary []Commentary
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, err2 := db.Prepare("SELECT * FROM Commentaire WHERE Post_ID = ?")
	if err2 != nil {
		fmt.Println(err2)
	}
	result, err3 := statement.Query(idPost)
	if err3 != nil {
		fmt.Println(err3)
	}
	for result.Next() {
		var Commentary Commentary
		var like string
		var dislike string
		result.Scan(&Commentary.IdCommentary, &Commentary.Text, &Commentary.IdPost, &Commentary.Username, &dislike, &like)
		Commentary.Dislike = KnowLike(dislike)
		Commentary.Like = KnowLike(like)
		Commentary.Username = GetUsernameByID(Commentary.Username)
		ListCommentary = append(ListCommentary, Commentary)
	}
	db.Close()
	return ListCommentary
}

func GetIdCategorie(categorieName string) string {
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, err2 := db.Prepare("SELECT ID_Categorie FROM Categorie WHERE Name = ?")
	if err2 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err2)
	}
	result, err3 := statement.Query(categorieName)
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

func GetUserPost(userID string) []Post {
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, err2 := db.Prepare("SELECT ID_Post,Text_Post,ID_Catégorie_Post, Like, Dislike FROM Post WHERE ID_User_Post = ?")
	if err2 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err2)
	}
	if err2 != nil {
		fmt.Println(err2)
	}
	result, err3 := statement.Query(userID)
	if err3 != nil {
		fmt.Println(err3)
	}
	var ListPost []Post
	for result.Next() {
		var singlePost Post
		var Like string
		var Dislike string
		var IdCategorie string
		result.Scan(&singlePost.IdPost, &singlePost.TextPost, &IdCategorie, &Like, &Dislike)
		singlePost.LikePost = KnowLike(Like)
		singlePost.DislikePost = KnowLike(Dislike)
		singlePost.Username = GetUsernameByID(userID)
		singlePost.CommentaryPost = GetCommmentary(singlePost.IdPost)
		singlePost.CategorieColor = GetColorCategoryById(IdCategorie)
		ListPost = append(ListPost, singlePost)
	}
	return ListPost
}

func GetUserPostLike(userID string) []Post {
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, err2 := db.Prepare("SELECT ID_User_Post,ID_Post,Text_Post,ID_Catégorie_Post, Like, Dislike FROM Post")
	if err2 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err2)
	}
	if err2 != nil {
		fmt.Println(err2)
	}
	result, err3 := statement.Query()
	if err3 != nil {
		fmt.Println(err3)
	}
	var ListPost []Post
	for result.Next() {
		var singlePost Post
		var Likestr string
		var Dislike string
		var IdCategorie string
		result.Scan(&singlePost.Username, &singlePost.IdPost, &singlePost.TextPost, &IdCategorie, &Likestr, &Dislike)
		if Like(Likestr, userID) {
			singlePost.LikePost = KnowLike(Likestr)
			singlePost.DislikePost = KnowLike(Dislike)
			singlePost.Username = GetUsernameByID(singlePost.Username)
			singlePost.CommentaryPost = GetCommmentary(singlePost.IdPost)
			singlePost.CategorieColor = GetColorCategoryById(IdCategorie)
			ListPost = append(ListPost, singlePost)
		}
	}
	return ListPost
}

func DeletePost(ID_Post string) {
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	prepareDeletePost, err2 := db.Prepare("DELETE FROM Post WHERE ID_Post = ?")
	if err2 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err2)
	}
	result, err3 := prepareDeletePost.Exec(ID_Post)
	if err3 != nil {
		fmt.Println("Erreur ouverture du fichier : ", err3)
	} else {
		fmt.Println(result)
	}
	db.Close()
}

func AddCategorie(Name string, color string) {
	db, err := sql.Open("sqlite3", "../BD/Forum_DB.db")
	if err != nil {
		fmt.Println("Erreur ouverture du fichier :")
		fmt.Println(err)
	}
	prepareAddCategorie, err := db.Prepare("INSERT INTO Categorie (Name, ID_Categorie, Color) VALUES (?,?,?)")
	var eRR error
	_, err2 := prepareAddCategorie.Exec(Name, uuid.Must(uuid.NewV4(), eRR), color)
	if err != nil || err2 != nil {
		fmt.Println("Erreur d'insertion :")
		fmt.Println(err)
		fmt.Println(err2)
	}
	db.Close()
}
