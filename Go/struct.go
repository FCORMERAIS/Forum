package main

type User struct {
	Id       string
	Username string
}

type Post struct {
	Username               string
	TextPost               string
	LikePost               int
	DislikePost            int
	IdPost                 string
	CommentaryPost         []Commentary
	CategorieColor         string
	CategorieName          string
	SamePersonWhithSession bool
}

type Categorie struct {
	Name  string
	Color string
}

type ForumPage struct {
	User           User
	ListCategories []Categorie
}

type Commentary struct {
	IdPost       string
	IdCommentary string
	Text         string
	Username     string
	Like         int
	Dislike      int
}
