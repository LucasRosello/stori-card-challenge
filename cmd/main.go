package main

import (
    "log"
    "github.com/LucasRosello/stori-card-challenge/pkg/database"
    "github.com/LucasRosello/stori-card-challenge/pkg/transactions"
    "github.com/LucasRosello/stori-card-challenge/pkg/mail"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    db, err := database.InitDB()
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer db.Close()


    trans, err := transactions.ReadFile("/root/transactions.csv")
    if err != nil {
        log.Fatal("Failed to read file:", err)
    }

 	insertTransactions(db, transactions)

    mail.SendNotificationMail(trans)
}