package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"encoding/csv"

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

	transactions := readFile()

	fmt.Println(transactions)
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

// readFile reads transactions from a CSV file and returns a slice of Transaction.
func readFile() []Transaction {
	file, err := os.Open("/root/transactions.csv")
	if err != nil {
		log.Fatalf("Unable to read input file: %s", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to parse file as CSV: %s", err)
	}

	var transactions []Transaction
	for i, record := range records {
		if i == 0 { // Skip the header row
			continue
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatalf("Error parsing ID: %s", err)
		}
		date := record[1]
		amount, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatalf("Error parsing transaction amount: %s", err)
		}
		transactions = append(transactions, Transaction{ID: id, Date: date, Amount: amount})
	}

	return transactions
}