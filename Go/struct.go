package main

type User struct {
	Id       string
	Username string
}

type Post struct {
	IDUser      string
	TextPost    string
	LikePost    string
	DislikePost string
}

type ArrayPosts struct {
	arrayPosts []Post
}
