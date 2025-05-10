package ctrl

import "net/http"

func Register(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Register endpoint"))
}

func Login(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Login endpoint"))
}