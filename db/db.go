// db/db.go
package db

import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the SQLite database
func InitDB() {
    var err error
    DB, err = sql.Open("sqlite3", "./linkshare.db")
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }

    // Create tables
    createTables := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS devices (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        name TEXT NOT NULL,
        last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );
    CREATE TABLE IF NOT EXISTS links (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        from_device_id INTEGER NOT NULL,
        to_device_id INTEGER NOT NULL,
        url TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (from_device_id) REFERENCES devices(id),
        FOREIGN KEY (to_device_id) REFERENCES devices(id)
    );
    `
    _, err = DB.Exec(createTables)
    if err != nil {
        log.Fatal("Failed to create tables:", err)
    }
}
