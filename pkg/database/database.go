package database

import (
    "database/sql"
    "log"
    "fmt"

    "github.com/LucasRosello/stori-card-challenge/pkg/transactions"
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

func InsertTransactions(db *sql.DB, transactions []transactions.Transaction) error {
    insertSQL := `INSERT INTO transactions (id, date, amount) VALUES (?, ?, ?)`
    stmt, err := db.Prepare(insertSQL)
    if err != nil {
        return fmt.Errorf("error preparing statement: %w", err)
    }
    defer stmt.Close()

    for _, trans := range transactions {
        _, err := stmt.Exec(trans.ID, trans.Date, trans.Amount)
        if err != nil {
            log.Printf("Error inserting transaction with ID %d: %v", trans.ID, err)
            continue
        }
    }

    log.Println("Transactions inserted successfully.")
    return nil
}