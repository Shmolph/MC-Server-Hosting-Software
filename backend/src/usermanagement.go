package main

import(
	"fmt"
	"database/sql"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"
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

func givemesessiontoken(db *sql.DB, userid int) (string, bool) {
	sessiontoken := rand.Text()
	tokenhashed := sha256.Sum256([]byte(sessiontoken))
	hashtext := fmt.Sprintf("%x", tokenhashed)
	if !storesession(db, hashtext, userid) {
		return "", false
	}
	return sessiontoken, true
}
