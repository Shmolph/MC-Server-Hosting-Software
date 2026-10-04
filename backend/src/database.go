package main

import(
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func connectdb() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=sam password=cats dbname=samsmctesting sslmode=disable")
	if err != nil {
		fmt.Println("DB: Failed to connect to the database!")
		fmt.Println(err)
	} else {
		fmt.Println("DB: DB handle created, testing connection....")
		pingstart := time.Now()
		err = db.Ping()
		if err != nil {
			fmt.Println("DB: Database:")
			fmt.Println(err)
		} else {
			fmt.Println("DB: DB Is connected with a Ping of", time.Since(pingstart))
		}
	}
	return db, err
}

func inittables(db *sql.DB) error {
	tables := []string{
		"CREATE TABLE IF NOT EXISTS userdata (username TEXT UNIQUE NOT NULL, password_hash TEXT NOT NULL, profile JSONB, id SERIAL PRIMARY KEY)",
		//"CREATE TABLE IF NOT EXISTS testtable (hello TEXT, world TEXT)",
	}
	for _, x := range tables {
		_, err := db.Exec(x)
		if err != nil {
			fmt.Println("DB: DATABASE: Fail in Initialising Tables:")
			fmt.Println("DB: %v", err)
			return err
		}
	}
	return nil
}
