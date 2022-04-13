package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func createConnection() *sql.DB {

	//Loading environment variables
	godotenv.Load()
	host := os.Getenv("HOST")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	//Database Connection String
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	//Opening connection to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// check the conection
	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Database Connection Successful!")
	return db
}

//Inserts new customer into the database
func insertCustomer(c Customer) int64 {
	//Connect to database
	db := createConnection()

	defer db.Close()

	//Create Customer
	sqlStatement := `
	INSERT INTO customer
	(firstname, lastname, email, phonenumber, occupation, customercity, createdat)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING customerid;`

	var id int64

	err := db.QueryRow(sqlStatement, c.FirstName, c.Lastname, c.Email, c.PhoneNumber, c.Occupation, c.CustomerCity, c.CreatedAt).Scan(&id)

	if err != nil {
		log.Printf("Unable to execute query. %v", err)
	}

	fmt.Printf("New Customer ID created - %v", id)

	//return the inserted id
	return id
}

//Initializes newly created customers account balance to 0
func initAccountBalance() {
	db := createConnection()

	defer db.Close()

	//Initialize account balance statement
	sqlStatement := `insert into accountbalance (accountnumber, customerid, accountbalance , balancedate)
	values (DEFAULT, DEFAULT, 0, current_timestamp) RETURNING accountbalance;`

	var acc_bal int64

	err := db.QueryRow(sqlStatement).Scan(&acc_bal)

	if err != nil {
		log.Print("Unable to initialize account", err)
	}
	fmt.Printf(", with Account Balance - %v", acc_bal)

}

//Updates a customers bank account after every transaction - deposits and withdrawals
func updateAccountBalance(id int64, accountbalance AccountBalance) int64 {
	db := createConnection()

	defer db.Close()

	//Initialize account balance statement
	sqlStatement := `update accountbalance set accountbalance =
	(select sum(CASE WHEN t.transactiontype = 'Credit' THEN t.transactionamount ELSE -t.transactionamount end)
	from transactionhistory t where accountnumber = $1
	), balancedate = now()
	where accountnumber =$1 RETURNING accountbalance;`

	var updated_acc_bal int64

	//Execute sql statement
	err := db.QueryRow(sqlStatement, id).Scan(&updated_acc_bal)
	if err != nil {
		log.Print("Unable to update account balance", err)
	}
	fmt.Printf(" - Your new Account Balance is - %v", updated_acc_bal)
	return updated_acc_bal
}

//Checks a specific customers bank account balance
func getBalance(id int64) (AccountBalance, error) {
	//Create the database connection
	db := createConnection()

	defer db.Close()

	//Create an accountbalance of AccountBalance type
	var accountbalance AccountBalance

	//Prepare account balance statement
	sqlStatement := `select * from accountbalance where customerid =$1;`

	//Execute sql statement
	row := db.QueryRow(sqlStatement, id)

	//Unmarshal the row object to accountbalance
	err := row.Scan(&accountbalance.AccountNumber, &accountbalance.CustomerID, &accountbalance.AccountBalance, &accountbalance.BalanceDate)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("This customer ID is non-existent!")
		return accountbalance, nil
	case nil:
		return accountbalance, nil
	default:
		log.Printf("Unable to scan row. %v", err)
	}
	return accountbalance, nil
}

//Makes Deposits into transaction history table
func depositTransactionHistory(id int64, transactionhistory TransactionHistory) int64 {
	//Create the database connection
	db := createConnection()

	defer db.Close()

	//Prepare account balance statement
	sqlStatement := `insert into transactionhistory (transactionnumber, accountnumber, transactiondate , currency, mediumoftransaction, transactiontype, transactionamount)
	values (DEFAULT, $1, now(), '$', $2, 'Credit', $3) RETURNING transactionamount;`

	//Store the transaction amount in a variable
	var trans_amount int64

	//Execute sql statement
	//Scan function will save the returned transaction amount
	err := db.QueryRow(sqlStatement, transactionhistory.AccountNumber, transactionhistory.MediumOfTransaction, transactionhistory.TransactionAmount).Scan(&trans_amount)
	if err != nil {
		log.Printf("Unable to complete transaction!. %v", err)
	}
	fmt.Printf(" - New deposit completed successfully, Amount - %v", trans_amount)

	// Return the transaction Amount
	return trans_amount
}

//Makes withdrawals - inserts entry into transaction history table
func withdrawalTransactionHistory(id int64, transactionhistory TransactionHistory) int64 {
	//Create the database connection
	db := createConnection()

	defer db.Close()

	//Prepare account balance statement
	sqlStatement := `insert into transactionhistory (transactionnumber, accountnumber, transactiondate , currency, mediumoftransaction, transactiontype, transactionamount)
	values (DEFAULT, $1, now(), '$', $2, 'Debit', $3) RETURNING transactionamount;`

	//Store the transaction amount in a variable
	var trans_amount int64

	//Execute sql statement
	//Scan function will save the returned transaction amount
	err := db.QueryRow(sqlStatement, transactionhistory.AccountNumber, transactionhistory.MediumOfTransaction, transactionhistory.TransactionAmount).Scan(&trans_amount)

	if err != nil {
		log.Printf("Unable to complete transaction!. %v", err)
	}
	fmt.Printf(" - New withdrawal completed successfully, Amount - %v", trans_amount)

	// Return the transaction Amount
	return trans_amount
}

//Handler function for performing the create customer task when a request is made to the create_customer route
func create_customer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//Create an empty customer of type Customer
	var c Customer

	//Decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		log.Print("Unable to decode the request body.", err)
	}

	// call inser customer function and pass the customer
	insertID := insertCustomer(c)
	initAccountBalance()

	// format a response object
	res := response{
		ID:       insertID,
		Message1: "Account created successfully, Thank you for using Mondu!",
	}

	//Send a response
	json.NewEncoder(w).Encode(res)

}

//Handler function for performing the deposit_cash task when a request is made to the deposit_cash route
func deposit_cash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//Get the CustomerID from the request params, key is "id"
	params := mux.Vars(r)

	//Convert the id type from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Unable to convert string into int", err)
	}
	//Create an empty transaction of type TransactionHistory
	var transactionhistory TransactionHistory
	var accountbalance AccountBalance

	//Decode the json request to transactionhistory
	err = json.NewDecoder(r.Body).Decode(&transactionhistory)
	if err != nil {
		log.Printf("Unable to decode the request body. %v", err)
	}

	//Call update updateTransactionHistory to perform transaction
	updatedRow1 := depositTransactionHistory(int64(id), transactionhistory)
	updatesRow2 := updateAccountBalance(int64(id), accountbalance)

	//format the meesage string
	msg1 := fmt.Sprintf("New deposit completed successfully, Amount -  %v", updatedRow1)
	msg2 := fmt.Sprintf(" - Your new account balance is - %v", updatesRow2)

	//format the response message
	res := response{
		ID:       int64(id),
		Message1: msg1,
		Message2: msg2,
	}
	json.NewEncoder(w).Encode(res)

}

//Handler function for performing the cash withdrawal task when a request is made to the withdraw_cash route
func withdraw_cash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//Get the CustomerID from the request params, key is "id"
	params := mux.Vars(r)

	//Convert the id type from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("Unable to convert string into int", err)
	}
	//Create an empty transaction of type TransactionHistory
	var transactionhistory TransactionHistory
	var accountbalance AccountBalance

	//Decode the json request to transactionhistory
	err = json.NewDecoder(r.Body).Decode(&transactionhistory)
	if err != nil {
		log.Printf("Unable to decode the request body. %v", err)
	}

	//Call Getbalance and ensure Customer has sufficient funds
	if transactionhistory.TransactionAmount >= accountbalance.AccountBalance {
		log.Println("Transaction declined - Insufficient funds.")
		json.NewEncoder(w).Encode("Insufficient Funds!")
	} else {
		//Call update updateTransactionHistory to perform transaction
		updatedRow1 := withdrawalTransactionHistory(int64(id), transactionhistory)
		updatedRow2 := updateAccountBalance(int64(id), accountbalance)

		//format the meesage string
		msg1 := fmt.Sprintf("New withdrawal completed successfully, Amount - %v", updatedRow1)
		msg2 := fmt.Sprintf(" - Your new account balance is - %v", updatedRow2)

		//format the response message
		res := response{
			ID:       int64(id),
			Message1: msg1,
			Message2: msg2,
		}
		json.NewEncoder(w).Encode(res)
	}

}

//Handler function for performing account balance task when a request is made to the account_balance route
func account_balance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//Get the CustomerID from the request params, key is "id"
	params := mux.Vars(r)

	//Convert the id type from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Print("Unable to convert string into int", err)
	}

	//Call the getBalance function with Customer ID to retrive single entry
	AccountBalance, err := getBalance(int64(id))
	if err != nil {
		log.Printf("Unable to get Account Balance. %v", err)
	}
	if AccountBalance.AccountNumber == "" {
		json.NewEncoder(w).Encode("Account Non-existent!")
	} else {
		//Send the response
		json.NewEncoder(w).Encode(AccountBalance)
	}

}
