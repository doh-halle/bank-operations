package main

import "time"

// Client Struct - Model
type Customer struct {
	CustomerID   string    `json:"customerid"`
	FirstName    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phonenumber"`
	Occupation   string    `json:"occupation"`
	CustomerCity string    `json:"customercity"`
	CreatedAt    time.Time `json:"createdate"`
}

type AccountBalance struct {
	AccountNumber  string    `json:"accountnumber"`
	CustomerID     string    `json:"customerid"`
	AccountBalance string    `json:"accountbalance"`
	BalanceDate    time.Time `json:"balancedate"`
}

type TransactionHistory struct {
	TransactionNumber   string    `json:"transactionnumber"`
	AccountNumber       string    `json:"accountnumber"`
	TransactionDate     time.Time `json:"transactiondate"`
	Currency            string    `json:"currency"`
	MediumOfTransaction string    `json:"mediumoftransaction"`
	TransactionType     string    `json:"transactiontype"`
	TransactionAmount   string    `json:"transactionamount"`
}

type response struct {
	ID       int64  `json:"id,omitempty"`
	Message1 string `json:"message1,omitempty"`
	Message2 string `json:"message2,omitempty"`
}
