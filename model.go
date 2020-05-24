package main

//Expense describe the main transaction object
type Expense struct {
	User        string `json:"user"`
	Date        string `json:"date"`
	Amount      string `json:"amount"`
	Description string `json:"description"`
	Category    string `json:"category"`
}
