package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"encoding/csv"
	"net/smtp"
    "bytes"
    "math"
    "html/template"
	"github.com/joho/godotenv"
    "strings"

	_ "github.com/mattn/go-sqlite3"
)

type Transaction struct {
	ID          int
	Date        string
	Amount		float64
}

type MonthTransactions struct {
    Month             string
    Transactions      []Transaction
    AverageDebit      float64
    AverageCredit     float64
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
	  
    monthTrans := []MonthTransactions{
        // TODO Add all the months
        {"Marzo", filterTransactions(transactions, "03"), 0, 0},
        {"Febrero", filterTransactions(transactions, "02"), 0, 0},
        {"Enero", filterTransactions(transactions, "01"), 0, 0},
    }

    for i, _ := range monthTrans {
        monthTrans[i].AverageDebit, monthTrans[i].AverageCredit = calculateAverages(monthTrans[i].Transactions)
    }

    tmpl := `{{range .}}
    <tr style="background-color: #f2f2f2; color: #333;">
        <td colspan="4" style="text-align: center; padding: 10px; font-weight: bold; background-color: #00d180; color: white;">{{.Month}}</td>
    </tr>
    {{range .Transactions}}
    <tr>
        <td style="padding: 8px; border-bottom: 1px solid #ddd;">{{.Date}}</td>
        <td style="padding: 8px; border-bottom: 1px solid #ddd;">Nombre del negocio</td>
        <td style="padding: 8px; border-bottom: 1px solid #ddd;">{{if lt .Amount 0.00}}<span style="color: #ff0000;">Debit</span>{{else}}<span style="color: #00d180;">Credit</span>{{end}}</td>
        <td style="padding: 8px; border-bottom: 1px solid #ddd;">${{printf "%.2f" .Amount}}</td>
    </tr>    
    {{end}}
    <tr class="average-row" style="font-weight: bold;">
        <td colspan="2" style="padding: 10px; background-color: #e8e8e8;">Promedio de transacciones con débito de {{.Month}}:</td>
        <td colspan="2" style="padding: 10px; background-color: #e8e8e8;">${{printf "%.2f" .AverageDebit}}</td>
    </tr>
    <tr class="average-row" style="font-weight: bold;">
        <td colspan="2" style="padding: 10px; background-color: #e8e8e8;">Promedio de transacciones con crédito de {{.Month}}:</td>
        <td colspan="2" style="padding: 10px; background-color: #e8e8e8;">${{printf "%.2f" .AverageCredit}}</td>
    </tr>
    {{end}}`
    

t, err := template.New("webpage").Parse(tmpl)
if err != nil {
panic(err)
}

var buf bytes.Buffer

err = t.Execute(&buf, monthTrans)
if err != nil {
panic(err)
}

outputHTML := buf.String()


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
                  <table class="transactions-table">
                    <tr>
                        <th>Fecha</th>
                        <th>Negocio</th>
                        <th>Tipo</th>
                        <th>Monto</th>
                    </tr>
                `+outputHTML+`
                </table>
              </div>
          </div>
      </body>
      </html>`
      
  
    message := ""
    for k, v := range header {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + htmlBody

  
	  err = smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, recipients, []byte(message))
	  if err != nil {
		  panic(err)
	  }
  
	  println("Email sended successfully!")
}

func filterTransactions(transactions []Transaction, month string) (filtered []Transaction) {
    for _, t := range transactions {
        if strings.Contains(t.Date, "-"+month+"-") {
            filtered = append(filtered, t)
        }
    }
    return
}

func calculateAverages(transactions []Transaction) (avgDebit, avgCredit float64) {
    var sumDebit, sumCredit float64
    var countDebit, countCredit int

    for _, t := range transactions {
        if t.Amount < 0 { // Debit if the amount is negative
            sumDebit += math.Abs(t.Amount)
            countDebit++
        } else if t.Amount > 0 { // Credit if the amount is positive
            sumCredit += t.Amount
            countCredit++
        }
    }

    if countDebit > 0 {
        avgDebit = sumDebit / float64(countDebit)
    }
    if countCredit > 0 {
        avgCredit = sumCredit / float64(countCredit)
    }
    return
}
