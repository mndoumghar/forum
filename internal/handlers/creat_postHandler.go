package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"forum/internal/auth"
	"forum/internal/db"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := auth.CheckSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	// Handle GET request: Render the create post form
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/creat_post.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// Handle POST request: Process form submission
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// if content == "" {
		// 	http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		// 	return
		// }

		// Temporary user ID (replace with actual session-based user ID)
		// kayna form wkayna Form L3adiya
		// 3lach khddmt b lform Bach njbd Ga3 l values Li Fihom same Name f input chechbox f Html
		// W r.form['status'] Kaththoum F slice
		// ama formValues katkhd ghir valus li drty lih CheckBox f html
		// hna bghina Ka3 element Dyal Categorie

		status := r.Form["status"]
		// hna Hwlt Slice l String Hint Maymknch Tsift l DAta base Slice f Vazlues khaso ikon string li howa TEXT Aw varchar nvarchr hado likaynin f database
		//
		statusStr := strings.Join(status, " ")
		fmt.Println(status)

		post_id := r.FormValue("post_id")
		content := r.FormValue("content")

		// in this Fucnction Chech If Session was Exist If Existed THen Save Your User_Id

		// Insert post into database
		_, err = db.DB.Exec("INSERT INTO category(user_id, post_id, status, content) VALUES(?, ?, ? , ?)", user_id, post_id,statusStr, content)
		if err != nil {
			http.Error(w, "Failed to create postss", http.StatusInternalServerError)
			return
		}

		// Redirect to posts page
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}

	// Handle other methods
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
