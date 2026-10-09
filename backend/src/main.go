package main

import(
	"fmt"
	"log"
	"os"
	"encoding/json"
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
	http.HandleFunc("/servers", withcors(requirelogin(db, listservers(db))))
	http.HandleFunc("/servers/create", withcors(requirelogin(db, createserver(db, versions))))
	http.HandleFunc("/servers/setup", withcors(requirelogin(db, setupserver(db))))
	http.HandleFunc("/servers/power", withcors(requirelogin(db, changeserverstate(db))))
	http.HandleFunc("/servers/command", withcors(requirelogin(db, sendcommandhandler(db))))
	http.HandleFunc("/servers/logs", withcors(requirelogin(db, getlogshandler(db))))
	http.HandleFunc("/servers/access", withcors(requirelogin(db, accesshandler(db))))
	http.HandleFunc("/servers/delete", withcors(requirelogin(db, deleteserverhandle(db))))
	http.HandleFunc("/versions", withcors(listversions(ordered)))
	http.HandleFunc("/logout", withcors(revokesession(db)))
	log.Fatal(http.ListenAndServe(":3030", nil))


	defer db.Close()
}

// ReadNames
func readnames(path string) []string {
	var output []struct {
		Name string	`json:"name"`
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return []string{}
	}
	
	err = json.Unmarshal(content, &output)
	if err != nil {
		return []string{}
	}
	names := []string{}
	for _, x := range output {
		names = append(names, x.Name)
	}
	return names
}
