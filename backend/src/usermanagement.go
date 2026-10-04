package main

import(
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"database/sql"
)

func createuser(db *sql.DB, username, password string) error {
	password_hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("USERMAN: Failed to Hash Password:", err)
		return err
	}
	err = adduser(db, username, string(password_hashed))
	if err != nil {
		fmt.Println("USERMAN: DB Failed to Insert Password", err)
		return err
	}
	return nil
}
