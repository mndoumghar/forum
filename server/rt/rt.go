package rt

import ( "net/http" ; "forum/server/ctrl")

func InitRoutes() {
    http.HandleFunc("/forum/register", ctrl.Register)
    http.HandleFunc("/forum/login", ctrl.Login)
}