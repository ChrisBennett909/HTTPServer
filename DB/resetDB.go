package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := sql.Open("sqlite3", "./credDB.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Generate bcrypt hash
	password := ""
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Insert new admin
	_, err = db.Exec(`
		INSERT INTO users (username, password, isAdmin, sessionToken)
		VALUES (?, ?, ?, ?)
	`, "Admin", string(hashed), 1, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Admin reset complete.")
}
