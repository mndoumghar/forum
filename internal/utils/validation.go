package utils

import ( "errors" ;"regexp" )

// ValidateRegistrationInput checks if the registration fields are valid.


func ValidateRegistrationInput(user, email, password string) error{
	
	if user == "" || email == "" || password == "" { return errors.New("email, username and password are required")}

	//validate email format
	E := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !E.MatchString(email) { return errors.New("invalid email format") }
	//validate username
	if len(user) < 3 { return errors.New("username must be at least 3 characters long") }
	if len(user) > 20 { return errors.New("username must be at most 20 characters long") }
	if !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(user) { return errors.New("username can only contain letters, numbers, dots, underscores, and hyphens") }
	//validate email
	if len(email) < 5 { return errors.New("email must be at least 5 characters long") }
	if len(email) > 50 { return errors.New("email must be at most 50 characters long") }
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) { return errors.New("invalid email format") }
	//validate password
	if len(password) < 8 { return errors.New("password must be at least 8 characters long") }
	if len(password) > 20 { return errors.New("password must be at most 20 characters long") }
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) || 
	   !regexp.MustCompile(`[a-z]`).MatchString(password) || 
	   !regexp.MustCompile(`\d`).MatchString(password) || // checks for at least one digit
	   !regexp.MustCompile(`[@$!%*?&]`).MatchString(password) { // checks for at least one special character
		return errors.New("password must contain an uppercase letter, a lowercase letter, a number, and a special character") 
	}



	// //PASSWD 
	// if password == "" { return errors.New("password is required") }
	// if len(password) < 8 { return errors.New("password must be at least 8 characters long")}
	// if len(password) > 20 { return errors.New("password must be at most 20 characters long")}

return nil
}