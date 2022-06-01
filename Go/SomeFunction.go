package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func KnowLike(listeLike string) int {
	var Liker = strings.Split(listeLike, "#")
	fmt.Println(Liker)
	if listeLike == "" {
		return 0
	}
	return len(Liker) - 1
}

func Like(listeLike string, userID string) bool {
	var Liker = strings.Split(listeLike, "#")
	for _, b := range Liker {
		if b == userID {
			return true
		}
	}
	return false
}

func trieMostLike(ListPost []Post) []Post {
	return ListPost
}

func trieLowLike(ListPost []Post) []Post {
	return ListPost
}
