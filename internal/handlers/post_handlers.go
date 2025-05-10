package handlers

import (
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register endpoint"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login endpoint"))
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Posts endpoint"))
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Post endpoint"))
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Comment endpoint"))
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Like endpoint"))
}