package main

import (
	"fmt"
	"log"
	"net/http"
	// Uncomment and use the necessary imports when needed
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

type App struct {
  DB *sql.DB
}

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


func connectToDB(connectionString string) *sql.DB {
  //Open the database connection
	db, err := sql.Open("postgres", connectionString )
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	
  log.Println("Reached DB")
  
  return db;
}

func (db *App) handlePost(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Received POST data securely!")

		insertQuery := `INSERT INTO your_table (column1, column2) VALUES ($1, $2)`
		// Execute the insert query
    _, err := db.DB.Exec(insertQuery, "value1", "value2")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Data inserted successfully")
    return
	}


func main() {
  // Loading env variables from .env files
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	if os.Getenv("API_KEY") == "" {
		log.Fatal("API_KEY environment variable is not set")
	}

	// Check if the environment variables are set
	if dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("One or more environment variables are not set!")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=timescaledb sslmode=disable", dbUser, dbPassword, dbName)
  db := connectToDB(connStr);
  defer db.Close()
  
  app := &App{DB: db}

	//Set up the HTTP mux
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the API!")
	})

	protectedPostHandler := http.HandlerFunc(app.handlePost)
	mux.Handle("/data", apiKeyMiddleware(protectedPostHandler))

	// Start server
	log.Println("Server running on :8443")
	log.Fatal(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", mux)) // Moved ListenAndServe after handler setup

}
