package whatsapp

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waLog "go.mau.fi/whatsmeow/util/log"
	"github.com/mdp/qrterminal/v3"

	_ "github.com/mattn/go-sqlite3"
)

func StartWhatsAppClient() *whatsmeow.Client {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:transactions.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatalf("Failed to create SQL store: %v", err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalf("Failed to get device store: %v", err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	SetupClient(client)

	return client
}

func SetupClient(client *whatsmeow.Client) {
	if client.Store.ID == nil {
		handleNewClient(client)
	} else {
		connectClient(client)
	}
}

func handleNewClient(client *whatsmeow.Client) {
	qrChan, _ := client.GetQRChannel(context.Background())
	err := client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect client: %v", err)
	}
	for evt := range qrChan {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			fmt.Println("Scan QR Code to login")
		} else {
			fmt.Printf("Login event: %s\n", evt.Event)
		}
	}
}

func connectClient(client *whatsmeow.Client) {
	err := client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect client: %v", err)
	}
	SendWhatsappMessage(client)
}

func SendWhatsappMessage(client *whatsmeow.Client) {
	user := os.Getenv("WHATSAPP_NUMBER")
	recipientJID := types.JID{User: user, Server: "s.whatsapp.net"}
    messageText := generateMessage()
	message := waProto.Message{Conversation: &messageText}
        
	response, err := client.SendMessage(context.Background(), recipientJID, &message)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
	fmt.Printf("Message sent with response: %v\n", response)
}

func WaitForInterrupt(client *whatsmeow.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
	fmt.Println("WhatsApp client disconnected successfully.")
}

func generateMessage() string {
	return `Stori

	Hola Lucas, te compartimos el resumen de tus stori cards. Ante cualquier incomveniente comunicate con nosotros.
	
	Tu balance es $20.33
	
	Tus últimas transacciones:
	
	Marzo:
	- 2024-03-20: Nombre del negocio - Credit: $150.00
	- 2024-03-15: Nombre del negocio - Debit: $-20.00
	Promedio de transacciones con débito de Marzo: $-20.00
	Promedio de transacciones con crédito de Marzo: $150.00
	
	Febrero:
	- 2024-02-18: Nombre del negocio - Credit: $200.00
	- 2024-02-10: Nombre del negocio - Debit: $-30.00
	Promedio de transacciones con débito de Febrero: $-30.00
	Promedio de transacciones con crédito de Febrero: $200.00
	
	Enero:
	- 2024-01-25: Nombre del negocio - Credit: $250.00
	- 2024-01-10: Nombre del negocio - Debit: $-45.00
	Promedio de transacciones con débito de Enero: $-45.00
	Promedio de transacciones con crédito de Enero: $250.00
	`
}