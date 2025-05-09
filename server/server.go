package main

import ( "log" ; "net/http" ;  "forum/server/rt")

func main() {
    rt.InitRoutes()

    log.Println("The server is running on port 8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Error starting server: ", err)
    }
}