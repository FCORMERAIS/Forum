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
	result3, err := db.Prepare("UPDATE Categorie SET Color = ? WHERE Name = \"Politique\"") // A CHANGER
	if err != nil {
		fmt.Println(err)
	}
	_, err4 := result3.Exec("#63542c")
	if err4 != nil {
		fmt.Println(err4)
	}
	db.Close()

}

//background-color: #0a79d4;
