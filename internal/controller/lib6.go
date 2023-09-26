package controller

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Define MySQL database connection parameters
	dbUser := "your_username"
	dbPass := "your_password"
	dbHost := "localhost" // Change to your MySQL server's host
	dbPort := 3306        // Change to your MySQL server's port
	dbName := "your_database_name"

	// Create a DSN (Data Source Name) for the MySQL connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL database!")

	// You can perform your authentication logic or other database operations here.
	// For example, you can execute SQL queries to check user credentials.

	// Example query to check if a user with a given username and password exists
	username := "user"
	password := "password"
	query := "SELECT COUNT(*) FROM users WHERE username=? AND password=?"
	var count int
	err = db.QueryRow(query, username, password).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		fmt.Println("Authentication successful")
	} else {
		fmt.Println("Authentication failed")
	}
}
