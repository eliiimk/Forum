package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    _ "github.com/mattn/go-sqlite3"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./user_auth.db")
    if err != nil {
        log.Fatal(err)
    }

    createTable := `CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        email TEXT NOT NULL
    );`

    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    hashedPassword, err := hashPassword(password)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    id := uuid.New().String()
    _, err = db.Exec("INSERT INTO users (id, username, password, email) VALUES (?, ?, ?, ?)", id, username, hashedPassword, email)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "User registered successfully!")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    var hashedPassword string
    err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        } else {
            http.Error(w, "Server error", http.StatusInternalServerError)
        }
        return
    }

    if !checkPasswordHash(password, hashedPassword) {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    fmt.Fprintf(w, "Login successful!")
}

func main() {
    initDB()
    defer db.Close()

    fs := http.FileServer(http.Dir("./Pages Principales"))
    http.Handle("/", fs)

    http.HandleFunc("/signup", signupHandler)
    http.HandleFunc("/login", loginHandler)

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
