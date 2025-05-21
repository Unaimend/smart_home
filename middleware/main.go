package main

import (
	"fmt"
	"log"
	"net/http"
  "io/ioutil"
	"time"
	// Uncomment and use the necessary imports when needed
	"database/sql"
  "encoding/json"
	_ "github.com/lib/pq"
	"os"
)

type App struct {
	DB *sql.DB
}
type ClimateData struct {
  Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	TimeStamp   time.Time `json:"timestamp"`
	Location string `json:"timestamp"`
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
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Reached DB")

	return db
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

func (db *App) handleClimate(w http.ResponseWriter, r *http.Request) {
  
  body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

  var climateData ClimateData

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body, &climateData)
	if err != nil {
		http.Error(w, "Server Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Output the received data (or handle it as needed)
	fmt.Printf("Received temperature: %d°C, humidity: %d%% at %04d-%02d-%02dT%02d:%02d:%02dZ\n at %s", climateData.Temperature, climateData.Humidity,
		climateData.TimeStamp.Year(),   // Year (4 digits)
		climateData.TimeStamp.Month(),  // Month (1-12)
		climateData.TimeStamp.Day(),    // Day (1-31)
		climateData.TimeStamp.Hour(),   // Hour (0-23)
		climateData.TimeStamp.Minute(), // Minute (0-59)
		climateData.TimeStamp.Second() // Second (0-59)
	) 


	insertQuery := `INSERT INTO temperature (timestamp, temperature, unit) VALUES ($1, $2, $3)`
	// Execute the insert query
	_, err = db.DB.Exec(insertQuery, climateData.TimeStamp, climateData.Temperature, "C")
	if err != nil {
		log.Fatal(err)
	}

	insertQuery = `INSERT INTO humidity (timestamp, humidity, unit) VALUES ($1, $2, $3)`
	// Execute the insert query
	_, err = db.DB.Exec(insertQuery, climateData.TimeStamp, climateData.Temperature, "C")
	if err != nil {
		log.Fatal(err)
	}
}

func (db *App) handleDefault(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the API!")
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
	db := connectToDB(connStr)
	defer db.Close()

	app := &App{DB: db}

	//Set up the HTTP mux
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.handleDefault)
	protectedPostHandler := http.HandlerFunc(app.handlePost)
	protectedClimateHandler := http.HandlerFunc(app.handleClimate)
	mux.Handle("/data", apiKeyMiddleware(protectedPostHandler))
	mux.Handle("/climate", apiKeyMiddleware(protectedClimateHandler))

	// Start server
	log.Println("Server running on :8443")
	log.Fatal(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", mux)) // Moved ListenAndServe after handler setup

}
