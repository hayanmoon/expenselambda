package main

//Expense describe the main transaction object
type Expense struct {
	Expenseid   string  `json:"expenseid"`
	Timestamp   int64   `json:"timestamp"`
	Amount      string  `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
}
