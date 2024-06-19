package db

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./user_auth.db")
    if err!= nil {
        log.Fatal(err)
    }

    createTable := `CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        email TEXT NOT NULL
    );`

    _, err = db.Exec(createTable)
    if err!= nil {
        log.Fatal(err)
    }
}
