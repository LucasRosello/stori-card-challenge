package transactions

import (
    "encoding/csv"
    "io"
    "os"
    "strconv"
    "strings"
)

func ReadFile(filePath string) ([]Transaction, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    csvReader := csv.NewReader(file)
    var transactions []Transaction

    if _, err := csvReader.Read(); err != nil {
        return nil, err
    }

    for {
        record, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        id, err := strconv.Atoi(record[0])
        if err != nil {
            return nil, err
        }
        amount, err := strconv.ParseFloat(record[2], 64)
        if err != nil {
            return nil, err
        }

        transactions = append(transactions, Transaction{
            ID:     id,
            Date:   record[1],
            Amount: amount,
        })
    }

    return transactions, nil
}

func FilterTransactions(transactions []Transaction, month string) []Transaction {
    var filtered []Transaction
    for _, t := range transactions {
        if strings.Contains(t.Date, "-"+month+"-") {
            filtered = append(filtered, t)
        }
    }
    return filtered
}

func CalculateAverages(transactions []Transaction) (avgDebit, avgCredit float64) {
    var sumDebit, sumCredit float64
    var countDebit, countCredit int

    for _, t := range transactions {
        if t.Amount < 0 {
            sumDebit += t.Amount
            countDebit++
        } else {
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
