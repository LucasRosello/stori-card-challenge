package main

import (
    "log"
    "github.com/LucasRosello/stori-card-challenge/pkg/database"
    "github.com/LucasRosello/stori-card-challenge/pkg/transactions"
    "github.com/LucasRosello/stori-card-challenge/pkg/mail"
    "github.com/LucasRosello/stori-card-challenge/pkg/whatsapp"
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

    err = database.InsertTransactions(db, trans)
    if err != nil {
        log.Fatal("Failed to insert transactions on DB", err)
    }

    mail.SendNotificationMail(trans)


    // Extra: Send resume to whatsapp
    client := whatsapp.StartWhatsAppClient()
    
    whatsapp.SendWhatsappMessage(client)
    
    whatsapp.WaitForInterrupt(client)
}