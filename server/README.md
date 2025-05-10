# Forum Application


This is a Go-based forum application that supports user authentication, dashboard access, and forum registration.


## Features

- User registration and login
- Secure session management
- SQLite database integration
- Modular route and handler structure


## Endpoints
- `/forum/register` - Register a new user
- `/forum/login` - Log in an existing user
- `/dashboard` - Access the user dashboard (requires authentication)


## Setup


1. Clone the repository:
   ```bash
   git clone https://github.com/mndoumghar/forum.git
   cd forum
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

4. Access the application at `http://localhost:8080`.