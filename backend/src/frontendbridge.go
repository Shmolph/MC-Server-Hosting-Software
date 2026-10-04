package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"database/sql"
	"strings"
	"net/url"
	"os"
	"time"
	"crypto/sha256"
)

//Signup
func registerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	type user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Captcha string `json:"captchaToken"`
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

	if len(u.Password) > 72 {
		http.Error(w, "password too long (max 72 characters)", http.StatusBadRequest)
		return
	}

	if !verifycaptcha(u.Captcha) {
		http.Error(w, "captcha failed", http.StatusForbidden)
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

//Login
func loginhandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	type user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Captcha string `json:"captchaToken"`
	}

	var u user
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		fmt.Println("Bridge: Bad request to Login:", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	u.Username = strings.TrimSpace(u.Username)
	if u.Username == "" {
		fmt.Println("Bridge: Whitespaced Password found on Login:", err)
		http.Error(w, "username cant be empty", http.StatusBadRequest)
		return
	}

	if u.Password == "" {
		fmt.Println("Bridge: Whitespaced Username found on Login:", err)
		http.Error(w, "password cant be empty", http.StatusBadRequest)
		return
	}
	if len(u.Password) > 72 {
		http.Error(w, "password too long (max 72 characters)", http.StatusBadRequest)
		return
	}

	if !verifycaptcha(u.Captcha) {
		http.Error(w, "captcha failed", http.StatusForbidden)
		return
	}

	userid, ok := logincheckdbside(db, u.Username, u.Password)
	if !ok {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
			
	token, ok2 := givemesessiontoken(db, userid)
	if !ok2 {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name: 	"session",
		Value:	token,
		Path:	"/",
		HttpOnly:	true,
		MaxAge:	7*24*60*60,
	})

	fmt.Fprintln(w, "Logged in!")
	return

}
}

//Account
func account(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "no cookie :(", http.StatusUnauthorized)
		return
	}
	ok, userid := checksession(db, (fmt.Sprintf("%x", sha256.Sum256([]byte(c.Value)))))
	if !ok {
		http.Error(w, "expired cookie xD", http.StatusUnauthorized)
		return
	}
	fmt.Fprintln(w, "logged in as user:", userid)
}
}

//Captcha
func verifycaptcha(token string) bool {
	token = strings.TrimSpace(token)
	if token == "" {
		return false
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{"secret": {os.Getenv("TURNSTILE_SECRET")}, "response": {token}})
	if err != nil {
		fmt.Println("Captcha request failed:", err)
		return false
	}
	defer resp.Body.Close()

	var r map[string]any
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return false
	}
	fmt.Println("captcha:", r)
	return r["success"] == true
}

//Cors
func withcors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

