package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

func main() {
	var err error
	Id := uuid.Must(uuid.NewV4(), err)
	fmt.Println(Id) // ID PERMETTANT DE SAVOIR QUI EST-CE A AJOUTER DANS LA ABSE DE DONNÉES APRES
	db, err := sql.Open("sqlite3", "./BD/Forum_DB.db")
	if err != nil {
		fmt.Println(err)
	}
	result, err := db.Prepare("INSERT INTO Categorie (Name, ID_Categorie) VALUES (?,?)") // A CHANGER
	if err != nil {
		fmt.Println(err)
	}
	_, err2 := result.Exec("Joyeux", Id)
	if err2 != nil {
		fmt.Println(err2)
	}
	db.Close()
}
