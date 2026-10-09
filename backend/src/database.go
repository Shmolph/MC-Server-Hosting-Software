package main

import(
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
	"golang.org/x/crypto/bcrypt"
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
		"CREATE TABLE IF NOT EXISTS sessions (token_hash TEXT PRIMARY KEY, user_id INTEGER NOT NULL REFERENCES userdata(id) ON DELETE CASCADE, expires_at TIMESTAMPTZ NOT NULL)",
		"CREATE TABLE IF NOT EXISTS servers (id SERIAL PRIMARY KEY, owner INTEGER NOT NULL REFERENCES userdata(id) ON DELETE CASCADE, name TEXT NOT NULL, type TEXT NOT NULL, version TEXT NOT NULL, ram_mb INTEGER NOT NULL, status TEXT NOT NULL DEFAULT 'offline', port INTEGER NOT NULL UNIQUE, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW())",
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

func adduser(db *sql.DB, username, passwordHash string) error {
	_, err := db.Exec("INSERT INTO userdata (username, password_hash) VALUES ($1, $2)", username, passwordHash)
	if err != nil {
		fmt.Println("DB: Failed to create new user:", err)
		return err
	}
	fmt.Println("DB: Created new user!")
	return nil
}

func logincheckdbside(db *sql.DB, username, password string) (int, bool) {
	var storedhash string
	var userid int
	err := db.QueryRow("SELECT id, password_hash FROM userdata WHERE username = $1", username).Scan(&userid, &storedhash)
	if err != nil {
		fmt.Println("Failed to Find Hash for user:", err)
		return 0, false
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedhash), []byte(password))
	if err != nil {
		return 0, false
	}
	return userid, true
}

func storesession(db *sql.DB, tokenhash string, userid int) bool {
	_, err := db.Exec("INSERT INTO sessions (token_hash, user_id, expires_at) VALUES ($1, $2, NOW() + INTERVAL '7 days')", tokenhash, userid)
	if err != nil {
		fmt.Println("Failed to save tokenhash:", err)
		return false
	}
	return true
}


func checksession(db *sql.DB, token_hash string) (bool, int) {
	var userid int
	err := db.QueryRow("SELECT user_id FROM sessions WHERE token_hash = $1 AND expires_at > NOW()", token_hash).Scan(&userid)
	if err != nil {
		return false, 0
	}
	return true, userid
}

func deletesession(db *sql.DB, hash string) bool {
	_, err := db.Exec("DELETE FROM sessions WHERE token_hash = $1", hash)
	if err != nil {
		return false
	}
	return true
}

func addserver(db *sql.DB, userid int, name string, servertype string, version string, ram int) (bool, int) {
	var serverid int
	err := db.QueryRow("INSERT INTO servers (owner, name, type, version, ram_mb, port) VALUES ($1, $2, $3, $4, $5, (SELECT COALESCE(MAX(port), 25564) + 1 FROM servers)) RETURNING id", userid, name, servertype, version, ram).Scan(&serverid)
	if err != nil {
		fmt.Println("DB: Failed to register server:", err)
		return false, 0
	}
	return true, serverid
}

// Fetch Server

// The struct
type serverdiagram struct {
	Id int `json:"id"`
	Name string	`json:"name"`
	Type string	`json:"type"`
	Version string	`json:"version"`
	Status string	`json:"status"`
	RamMB int	`json:"ram_mb"`
	Port int	`json:"port"`
}

// Actual Function
func fetchserver(db *sql.DB, userid int) ([]serverdiagram, error) {
	rows, err := db.Query("SELECT id, name, type, version, status, ram_mb, port FROM servers WHERE owner = $1", userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	servers := []serverdiagram{}
	for rows.Next() {
		var s serverdiagram
		err := rows.Scan(&s.Id, &s.Name, &s.Type, &s.Version, &s.Status, &s.RamMB, &s.Port)
		if err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return servers, nil
}

// Get Server
func getserver(db *sql.DB, serverid int, userid int) (serverdiagram, bool) {
	var s serverdiagram
	err := db.QueryRow("SELECT id, name, type, version, status, ram_mb, port FROM servers WHERE id = $1 AND owner = $2", serverid, userid).Scan(&s.Id, &s.Name, &s.Type, &s.Version, &s.Status, &s.RamMB, &s.Port)
	if err != nil {
		return serverdiagram{}, false
	}
	return s, true
}

// Remove Server
func deleteserver(db *sql.DB, serverid int, userid int) bool {
	result, err := db.Exec("DELETE FROM servers WHERE id = $1 AND owner = $2", serverid, userid)
	if err != nil {
		fmt.Println("Failed to delete:", err)
		return false
	}
	rows, err := result.RowsAffected()
	if rows == 0 || err != nil {
		return false
	}
	return true
}
