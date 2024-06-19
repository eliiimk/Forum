package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"

    "Forum/backend/db" // import db package
)

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method!= http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    hashedPassword, err := hashPassword(password)
    if err!= nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    id := uuid.New().String()
    _, err = db.db.Exec("INSERT INTO users (id, username, password, email) VALUES (?,?,?,?)", id, username, hashedPassword, email)
    if err!= nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/login", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method!= http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    var hashedPassword string
    err := db.db.QueryRow("SELECT password FROM users WHERE username =?", username).Scan(&hashedPassword)
    if err!= nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        } else {
            http.Error(w, "Server error", http.StatusInternalServerError)
        }
        return
    }

    if!checkPasswordHash(password, hashedPassword) {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    http.Redirect(w, r, "/welcome", http.StatusFound)
}

func main() {
    db.initDB()
    defer db.db.Close()

    fs := http.FileServer(http.Dir("."))
    http.Handle("/", fs)

    http.HandleFunc("/signup", signupHandler)
    http.HandleFunc("/login", loginHandler)

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}