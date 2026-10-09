package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"encoding/json"
	"database/sql"
	"os"
	"strconv"
)

// Server Json
type minecraftversion struct {
	ID	string	`json:"id"`
	Type	string	`json:"type"`
}

//Create Server
func createserver(db *sql.DB, versions map[string]bool) (func(http.ResponseWriter, *http.Request, int)) {
	// Server Types
	types := []string{
		"Vanilla",
		"Forge",
		"Fabric",
		"Paper",
	}
	return func(w http.ResponseWriter, r *http.Request, userid int) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		type ingredients struct {
			Name string `json:"name"`
			Type string	`json:"type"`
			Version string	`json:"version"`
			Ram int	`json:"ram_mb"`
			Eula bool	`json:"eula"`
		}
		var i ingredients
		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			fmt.Println("Failed to decode json payload:", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		i.Name = strings.TrimSpace(i.Name)
		i.Type = strings.TrimSpace(i.Type)
		i.Version = strings.TrimSpace(i.Version)

		if i.Name == "" {
			fmt.Println("SM: Bad request to Login:", err)
			http.Error(w, "name cant be empty", http.StatusBadRequest)
			return
		}
		if i.Type == "" {
			fmt.Println("SM: Bad request to Login:", err)
			http.Error(w, "type cant be empty", http.StatusBadRequest)
			return
		}
		if i.Version == "" {
			fmt.Println("SM: Bad request to Login:", err)
			http.Error(w, "version cant be empty", http.StatusBadRequest)
			return
		}

		if !versions[i.Version] {
			http.Error(w, "version dosent exist! (skill issue?)", http.StatusBadRequest)
			return
		}

		if len(i.Name) > 32 {
			http.Error(w, "name too long (max 32 characters)", http.StatusBadRequest)
			return
		}

		if len(i.Name) < 3 { // <3 lol
			http.Error(w, "name too short (min 3 characters)", http.StatusBadRequest)
			return
		}

		if i.Ram < 1024 {
			http.Error(w, "ram must be 1024 - 4096", http.StatusBadRequest)
			return
		}

		if i.Ram > 4096 {
			http.Error(w, "ram must be 1024 - 4096", http.StatusBadRequest)
			return
		}

		if !i.Eula {
			http.Error(w, "eula must be true!", http.StatusBadRequest)
			return
		}

		correct := false
		for _, x := range types {
			if x== i.Type {
				correct = true
			}
		}
		if !correct {
			http.Error(w, "type must be, 'Vanilla', 'Forge', 'Fabric', Or 'Paper'", http.StatusBadRequest)
			return
		}

		ok, serverid := addserver(db, userid, i.Name, i.Type, i.Version, i.Ram)
		if !ok {
			http.Error(w, "failed to create server", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, serverid)
	}
}

// List Servers
func listservers(db *sql.DB) func(http.ResponseWriter, *http.Request, int) {
	return func(w http.ResponseWriter, r *http.Request, userid int) {
		servers, err := fetchserver(db, userid)
		if err != nil {
			http.Error(w, "Failed to fetch server!", http.StatusInternalServerError)
			return
		}
		for i := range servers {
			servers[i].Status = servercontainerstatus(servers[i].Id)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(servers)
	}
}

// Server Type Manifest
func pullmanifest() (map[string]bool, []string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://piston-meta.mojang.com/mc/game/version_manifest_v2.json")
	if err != nil {
		fmt.Println("SM: Failed to connect to mojang version manifest!")
		return nil, nil, err
	}
	defer resp.Body.Close()
	type mcmanifest struct {
		Versions []minecraftversion	`json:"versions"`
	}
	var m mcmanifest
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, nil, err
	}
	versions := make(map[string]bool)
	ordered := []string{}
	for _, x := range m.Versions {
		if x.Type == "release" {
			versions[x.ID] = true
			ordered = append(ordered, x.ID)
		}
	}

	return versions, ordered, nil	
}

// List Versions
func listversions(ordered []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ordered)
	}
}

// Setup Server
func setupserver(db *sql.DB) func(http.ResponseWriter, *http.Request, int) {
	return func(w http.ResponseWriter, r *http.Request, userid int) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		type idholder struct {
			Id int	`json:"id"`
		}
		var h idholder
		err := json.NewDecoder(r.Body).Decode(&h)
		if err != nil {
			http.Error(w, "please send valid json!", http.StatusBadRequest)
			return
		}
		server, ok := getserver(db, h.Id, userid)
		if !ok {
			http.Error(w, "failed to find server", http.StatusNotFound)
			return
		}
		err = createcontainer(server)
		if err != nil {
			http.Error(w, "failed to create server", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "server setup")
		return
	}
}

// Boschinate Server
func changeserverstate(db *sql.DB) func(http.ResponseWriter, *http.Request, int) {
	return func(w http.ResponseWriter, r *http.Request, userid int) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		type holder struct {
			Id int	`json:"id"`
			Action string	`json:"action"`
		}
		var h holder
		err := json.NewDecoder(r.Body).Decode(&h)
		if err != nil {
			http.Error(w, "please send valid json T_T", http.StatusBadRequest)
			return
		}
		server, ok := getserver(db, h.Id, userid)
		if !ok {
			http.Error(w, "you dont own this server", http.StatusNotFound)
			return
		}
		switch h.Action {
		case "start":
			err = changestatus(server.Id, "start")
		case "stop":
			err = changestatus(server.Id, "stop")
		case "restart":
			err = changestatus(server.Id, "restart")
		default:
			http.Error(w, "unknown action", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "ok")
	}
}

// Delete Server
func deleteserverhandle(db *sql.DB) func(http.ResponseWriter, *http.Request, int) {
	return func(w http.ResponseWriter, r *http.Request, userid int) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		type holder struct {
			Id int	`json:"id"`
		}
		var h holder
		err := json.NewDecoder(r.Body).Decode(&h)
		if err != nil {
			http.Error(w, "failed to decode json, please send proper json", http.StatusBadRequest)
			return
		}
		server, ok := getserver(db, h.Id, userid)
		if !ok {
			http.Error(w, "failed to find server", http.StatusNotFound)
			return
		}
		switch servercontainerstatus(server.Id) {
		case "starting", "online":
			http.Error(w, "please stop server before deleting", http.StatusConflict)
			return
		case "offline":
			err := changestatus(server.Id, "rm")
			if err != nil {
				http.Error(w, "failed to change server status", http.StatusInternalServerError)
				return
			}
		case "not setup":

		}
		if !deleteserver(db, server.Id, userid) {
			http.Error(w, "failed to delete server", http.StatusInternalServerError)
			return
		}
		err = os.RemoveAll("/home/mc-servers/" + strconv.Itoa(server.Id))
		if err != nil {
			fmt.Println("SERVERMAN:", err)
		}
		fmt.Fprintln(w, "ok")
	}
}

// Send command
func sendcommandhandler(db *sql.DB) func(http.ResponseWriter, *http.Request, int) {
	return func(w http.ResponseWriter, r *http.Request, userid int) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		type holder struct {
			Id int	`json:"id"`
			Command string	`json:"command"`
		}
		var h holder
		err := json.NewDecoder(r.Body).Decode(&h)
		if err != nil {
			http.Error(w, "Failed to decode json, please send proper json", http.StatusBadRequest)
			return
		}
		server, ok := getserver(db, h.Id, userid)
		if !ok {
			http.Error(w, "Failed to find server", http.StatusNotFound)
			return
		}
		if servercontainerstatus(server.Id) != "online" {
			http.Error(w, "Please wait for server to be started before sending a command!", http.StatusConflict)
			return
		}
		cmd := strings.TrimSpace(h.Command)
		if cmd == "" {
			http.Error(w, "Bad command", http.StatusBadRequest)
			return
		}
		if len(cmd) > 256 {
			http.Error(w, "Bad command", http.StatusBadRequest)
			return
		}
		if strings.ContainsAny(cmd, "\r\n") {
			http.Error(w, "Bad command", http.StatusBadRequest)
			return
		}
		if strings.HasPrefix(cmd, "-") {
			http.Error(w, "Bad command", http.StatusBadRequest)
			return
		}
		output, err := sendcommand(server.Id, cmd)
		if err != nil {
			http.Error(w, output, http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, output)
	}
}

// Get Logs
func getlogshandler(db *sql.DB) func(http.ResponseWriter, *http.Request, int) {
	return func(w http.ResponseWriter, r *http.Request, userid int) {
		idstr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			http.Error(w, "bad request (failed to convert)", http.StatusBadRequest)
			return
		}
		server, ok := getserver(db, id, userid)
		if !ok {
			http.Error(w, "server not found", http.StatusNotFound)
			return
		}
		if servercontainerstatus(server.Id) == "not setup" {
			http.Error(w, "server is not setup", http.StatusConflict)
			return
		}
		output, err := getlogs(server.Id) 
		if err != nil {
			http.Error(w, "failed to get logs", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, output)
	}
}

// Access Handler
func accesshandler(db *sql.DB) func(http.ResponseWriter, *http.Request, int) {
	return func(w http.ResponseWriter, r *http.Request, userid int) {
		idstr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		server, ok := getserver(db, id, userid)
		if !ok {
			http.Error(w, "failed to find server", http.StatusNotFound)
			return
		}
		dir := "/home/mc-servers/" + strconv.Itoa(server.Id)
		type response struct {
			Ops []string	`json:"ops"`
			Banned []string	`json:"banned"`
		}
		resp := response{
			Ops:	readnames(dir + "/ops.json"),
			Banned: readnames(dir + "/banned-players.json"),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
