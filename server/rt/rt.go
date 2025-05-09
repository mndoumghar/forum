package rt

import ( "net/http" ; "fmt" ; "os")

func InitRoutes() {
    // serve static files from the "static" directory
    http.Handle("/", http.FileServer(http.Dir("./static")))

    // server register and login pages
    http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            http.ServeFile(w, r, "static/register.html")
        } else if r.Method == http.MethodPost {
            // handle registration logic here
            // get form values
            username := r.FormValue("username")
            password := r.FormValue("password")
            email := r.FormValue("email")
            // check if username already exists in the database
            // if it does, return an error
            // if not, hash the password and store the user in the database
            // use bcrypt to hash the password
            // hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
            // Open the file in append mode, create it if it doesn't exist
            file, err := os.OpenFile("server/DataBase/data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }
            defer file.Close()

            // Write the user data to the file
            _, err = file.WriteString(fmt.Sprintf("Username: %s, Email: %s, Password: %s\n", username, email, password))
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            // Redirect to success page
            http.Redirect(w, r, "/success", http.StatusSeeOther)
            //redirect to success page
        }
    })

    // http.HandleFunc("/forum/register", ctrl.Register)
    // http.HandleFunc("/forum/login", ctrl.Login)
}