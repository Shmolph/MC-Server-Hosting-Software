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



	http.HandleFunc("/register", registerHandler(db))
	log.Fatal(http.ListenAndServe(":3030", nil))



	defer db.Close()
}
