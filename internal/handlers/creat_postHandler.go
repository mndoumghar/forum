package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"forum/internal/db"
)

func CreatpostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(db.GetPost())
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/creat_post.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Error Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

}
