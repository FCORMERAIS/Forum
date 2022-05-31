package main

type User struct {
	Id       string
	Username string
}

type Post struct {
	IDUser      string
	TextPost    string
	LikePost    int
	DislikePost int
}

type ArrayPosts struct {
	arrayPosts []Post
}

type Categorie struct {
	URL  string
	Name string
}

type ForumPage struct {
	User           User
	ListCategories []Categorie
}
