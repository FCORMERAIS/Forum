package main

type User struct {
	Id       string
	Username string
}

type Post struct {
	IDUser      int
	TextPost    string
	LikePost    bool
	DislikePost bool
}

type ArrayPosts struct {
	arrayPosts []Post
}
