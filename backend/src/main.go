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
	
	fmt.Println("Pulling Version Manifest,,,")
	versions, ordered, err := pullmanifest()
	if err != nil {
		log.Fatalf("Failed to pull manifest: %v", err)
	}
	fmt.Println(len(versions), "versions loaded,", len(ordered), "in list")

	// Networking
	http.HandleFunc("/register", withcors(registerHandler(db)))
	http.HandleFunc("/login", withcors(loginhandler(db)))
	http.HandleFunc("/account", withcors(requirelogin(db, account)))
	http.HandleFunc("/servers/create", withcors(requirelogin(db, createserver(db, versions))))
	http.HandleFunc("/servers", withcors(requirelogin(db, listservers(db))))
	http.HandleFunc("/versions", withcors(listversions(ordered)))
	http.HandleFunc("/logout", withcors(revokesession(db)))
	log.Fatal(http.ListenAndServe(":3030", nil))





	defer db.Close()
}
