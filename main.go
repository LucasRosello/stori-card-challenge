package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Transaction struct {
	ID          int
	Date        string
	Amount		float64
}

func main() {
	db, err := sql.Open("sqlite3", "/root/db/transactions.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db)

}

// createTable creates a table in the database if it does not exist.
func createTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS transactions (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		"date" TEXT, 
		"amount" REAL
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	}

	fmt.Println("Database successfully working.")
}