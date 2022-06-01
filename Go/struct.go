package main

type User struct {
	Id       string
	Username string
}

type Post struct {
	Username    string
	TextPost    string
	LikePost    int
	DislikePost int
	IdPost      string
}

type Categorie struct {
	URL  string
	Name string
}

type ForumPage struct {
	User           User
	ListCategories []Categorie
}
