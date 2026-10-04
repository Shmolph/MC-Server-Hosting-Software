package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"database/sql"
	"strings"
)

func registerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	type user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var u user
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		fmt.Println("Bridge: Oh noes couldnt decode json:", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	u.Username = strings.TrimSpace(u.Username)

	if u.Username == "" {
		http.Error(w, "username cant be empty", http.StatusBadRequest)
		return
	}

	if u.Password == "" {
		http.Error(w, "password cant be empty", http.StatusBadRequest)
		return
	}

	err = createuser(db, u.Username, u.Password)
	if err != nil && !strings.Contains(err.Error(), "23505") {
		fmt.Println("Bridge: Failed to add user", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	} else if err != nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}
	fmt.Fprintln(w, "Received!")
}
}