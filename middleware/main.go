package main

import (
	"fmt"
	"log"
	"net/http"
	// Uncomment and use the necessary imports when needed
	// "database/sql"
	 "os"
	// _ "github.com/lib/pq"
)


var validAPIKeys = map[string]bool{
    os.Getenv("API_KEY"): true, // API key loaded from environment variable
}

// Middleware for API key authentication
func apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		 if !validAPIKeys[apiKey] {
			log.Println(apiKey)
		 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		 	return
		 }
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Uncomment and replace the database connection setup when needed
	// dbUser := os.Getenv("POSTGRES_USER")
	// dbPassword := os.Getenv("POSTGRES_PASSWORD")
	// dbName := os.Getenv("POSTGRES_DB")
 	if os.Getenv("API_KEY") == "" {
        	log.Fatal("API_KEY environment variable is not set")
    	}

	// Check if the environment variables are set
	// if dbUser == "" || dbPassword == "" || dbName == "" {
	// 	log.Fatal("One or more environment variables are not set!")
	// }

	// connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=timescaledb sslmode=disable", dbUser, dbPassword, dbName)
	// Open the database connection
	// log.Println("Trying to open DB 1")	
	// db, err := sql.Open("postgres", connStr)
	// log.Println("Trying to open DB 2")	
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// Check if the database is reachable
	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Reached DB")	

	// Set up the HTTP mux
	mux := http.NewServeMux()

	// Public route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the API!")
	})

	// Protected route (requires API key)
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "You accessed a protected route!")
	})
	mux.Handle("/protected", apiKeyMiddleware(protectedHandler)) // Fixed route handler registration

	// Protected POST route
	protectedPostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Received POST data securely!")
	})
	mux.Handle("/data", apiKeyMiddleware(protectedPostHandler)) // Fixed route handler registration

	// Start server
	log.Println("Server running on :8443")
	log.Fatal(http.ListenAndServe(":8443", mux)) // Moved ListenAndServe after handler setup
}

