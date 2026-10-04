package main

import(
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("MAIN: Initialising DB")
	db, err := connectdb()
	if err != nil {
		log.Fatalf("MAIN: UNABLE TO USE DB: %v", err)
	} else {
		fmt.Println("MAIN: DB Working Fine From Here!")
	}
	err = inittables(db)
	if err != nil {
		log.Fatalf("MAIN: UNABLE TO CREATE TABLES: %v", err)
	} else {
		fmt.Println("MAIN: If Tables Missing, they are created")
	}
	go func() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("/home/shmolph/Documents/Personal/Fun Things =D/MC Server Hosting Software/frontend/")))
	log.Fatal(http.ListenAndServe(":5000", mux))
	}()

	http.HandleFunc("/register", withcors(registerHandler(db)))
	http.HandleFunc("/login", withcors(loginhandler(db)))
	http.HandleFunc("/account", withcors(account(db)))
	log.Fatal(http.ListenAndServe(":3030", nil))





	defer db.Close()
}
