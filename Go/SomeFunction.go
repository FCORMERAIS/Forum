package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// fonction vérifiant si le mot de passe est assez sécurisé (10 caractère, 1 maj , 1 min, 1 caractère spécial, 1 chiffre)
func passwordGood(mdp string, w http.ResponseWriter) bool {
	if len(mdp) < 10 {
		fmt.Fprintf(w, `<p class="error_message">le MDP EST TROP COURT</p>`)
		return false
	}
	r, _ := regexp.Compile("1|2|3|4|5|6|7|8|9|0")
	if !r.MatchString(mdp) {
		fmt.Fprintf(w, `<p class="error_message">IL MANQUE UN CHIFFRE DANS VOTRE MDP</p>`)
		return false
	}
	r2, _ := regexp.Compile("/|#|;|!|$|}")
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

//fonction permettant de hasher un password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

//fonction permettant de vérifier si un password hasher est bien égal au mot de passe que rentre l'utilisateur
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//fonction permettant de savoir le nombre de like dans un post ou un commentaires grace a une chaine de caractère
func KnowLike(listeLike string) int {
	var Liker = strings.Split(listeLike, "#")
	if listeLike == "" {
		return 0
	}
	return len(Liker) - 1
}

//fonction permettant de vérifier si l'utilisateur a liker ou nan un post ou un commentaires grace a son id et a une string (renvoie un bool)
func Like(listeLike string, userID string) bool {
	var Liker = strings.Split(listeLike, "#")
	for _, b := range Liker {
		if b == userID {
			return true
		}
	}
	return false
}
