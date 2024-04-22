package transactions

type Transaction struct {
    ID     int
    Date   string
    Amount float64
}

type MonthTransactions struct {
    Month           string
    Transactions    []Transaction
    AverageDebit    float64
    AverageCredit   float64
}