package mail

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "log"
    "net/smtp"
    "os"
    "strings"
    "text/template"

    "github.com/LucasRosello/stori-card-challenge/pkg/transactions"
)

func SendNotificationMail(trans []transactions.Transaction) {
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
    senderEmail := os.Getenv("SMTP_USER")
    password := os.Getenv("SMTP_PASSWORD")
    recipients := []string{senderEmail}

    auth := smtp.PlainAuth("", senderEmail, password, smtpHost)
    
    htmlBody := generateHTMLForMail(trans)

    header := buildHeader(senderEmail, recipients)
    message := formatMessage(header, htmlBody)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, recipients, []byte(message))
    if err != nil {
        log.Fatalf("Failed to send email: %s", err)
    }

    log.Println("Email sent successfully!")
}

func loadTemplate() *template.Template {
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
    {{end}}`

    t, err := template.New("emailTemplate").Parse(tmpl)
    if err != nil {
        log.Fatalf("Failed to parse template: %s", err)
    }
    return t
}

func generateHTMLForMail(trans []transactions.Transaction) string {
    tmpl := loadTemplate()

    monthTrans := []transactions.MonthTransactions{
        {"Marzo", transactions.FilterTransactions(trans, "03"), 0, 0},
        {"Febrero", transactions.FilterTransactions(trans, "02"), 0, 0},
        {"Enero", transactions.FilterTransactions(trans, "01"), 0, 0},
    }

    outputHTML, err := executeTemplate(tmpl, monthTrans)
    if err != nil {
        log.Fatalf("Failed to execute template: %s", err)
    }

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
                <table class="transactions-table" style="width: 100%;">
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

    return htmlBody
}

func executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
    var buf bytes.Buffer
    err := tmpl.Execute(&buf, data)
    if err != nil {
        return "", err
    }
    return buf.String(), nil
}

func buildHeader(from string, to []string) map[string]string {
    header := make(map[string]string)
    header["From"] = from
    header["To"] = strings.Join(to, ",")
    header["Subject"] = "StoriCard Challenge"
    header["MIME-Version"] = "1.0"
    header["Content-Type"] = "text/html; charset=\"UTF-8\""
    header["Content-Transfer-Encoding"] = "base64"

    return header
}

func formatMessage(header map[string]string, body string) string {
    var message strings.Builder
    for k, v := range header {
        message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
    }
    encodedBody := base64.StdEncoding.EncodeToString([]byte(body))
    message.WriteString("\r\n" + encodedBody)

    return message.String()
}
