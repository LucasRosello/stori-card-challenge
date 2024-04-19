package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"encoding/csv"
	"net/smtp"
	"github.com/joho/godotenv"
    "strings"

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

	err = godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

	createTable(db)

	transactions := readFile()

	insertTransactions(db, transactions)

	SendNotificationMail(transactions)	
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

// insertTransactions inserts a slice of Transaction into the database.
func insertTransactions(db *sql.DB, transactions []Transaction) {
	for _, trans := range transactions {
		insertSQL := `INSERT INTO transactions (id, date, amount) VALUES (?, ?, ?)`
		statement, err := db.Prepare(insertSQL)
		if err != nil {
			log.Fatalf("Error preparing statement: %s", err)
			continue
		}

		_, err = statement.Exec(trans.ID, trans.Date, trans.Amount)
		if err != nil {
			log.Fatalf("Error inserting transaction data: %s", err)
			continue
		}
		
		fmt.Println("Transactions inserted successfully.")
	}
}


// SendNotificationMail send the resumee via mail using Gmail SMTP
func SendNotificationMail(transactions []Transaction) {
	  
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
    recipients := []string{senderEmail}
  
	  auth := smtp.PlainAuth("", senderEmail, password, smtpHost)
  
	  header := make(map[string]string)
	  header["From"] = senderEmail
	  header["To"] = strings.Join(recipients, ",")
	  header["Subject"] = "StoriCard Challenge"
	  header["MIME-Version"] = "1.0"
	  header["Content-Type"] = "text/html; charset=\"UTF-8\""
	  header["Content-Transfer-Encoding"] = "base64"
  
      htmlBody := `
      <!DOCTYPE html>
      <html>
      <head>
      <title>Resumen de Transacciones</title>
      </head>
      <body style="font-family: Arial, sans-serif; background-color: #f2f2f2; color: #000000; padding: 20px; margin: 0;">
          <div style="background-color: #ffffff; width: 100%; max-width: 600px; margin: 0 auto; padding: 20px; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
              <div style="background-color: #00d180; color: #ffffff; padding: 10px; text-align: center; border-radius: 8px 8px 0 0; font-size: 24px; font-weight: bold;">Stori</div>
              <div style="padding: 20px; text-align: center;">
                  <p>Hola <span style="color: #00d180;">Lucas</span>, te compartimos el resumen de tus stori cards. Ante cualquier incomveniente comunicate con nosotros.</p>
                  <div style="font-size: 22px; margin: 20px 0; color: #000000;">Tu balance es $20.33</div>
                  <div style="font-size: 22px; text-align: center; color: #003a40; margin: 20px 0;">Tus últimas transacciones</div>
                  <!-- El resto de tu HTML con estilos en línea -->
              </div>
          </div>
      </body>
      </html>`
      
  
message := ""
for k, v := range header {
	message += fmt.Sprintf("%s: %s\r\n", k, v)
}
message += "\r\n" + htmlBody

  
	  // Conectar al servidor SMTP y enviar el mensaje
	  err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, recipients, []byte(message))
	  if err != nil {
		  panic(err)
	  }
  
	  println("Email sended successfully!")
}