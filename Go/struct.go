package main

type User struct {
	Id       string
	Username string
}

type Post struct {
	Username       string
	TextPost       string
	LikePost       int
	DislikePost    int
	IdPost         string
	CategorieColor string
	CategorieName  string
}

type Categorie struct {
	URL   string
	Name  string
	Color string
}

type ForumPage struct {
	User           User
	ListCategories []Categorie
}
