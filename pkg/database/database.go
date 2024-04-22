package database

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "/root/db/transactions.db")
    if err != nil {
        return nil, err
    }

    if err := createTable(db); err != nil {
        return nil, err
    }

    return db, nil
}

func createTable(db *sql.DB) error {
    createTableSQL := `CREATE TABLE IF NOT EXISTS transactions (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
        "date" TEXT, 
        "amount" REAL
    );`
    _, err := db.Exec(createTableSQL)
    if err != nil {
        return err
    }

    log.Println("Database successfully configured.")
    return nil
}
