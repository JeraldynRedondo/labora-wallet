package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// PostgresDBHandler contains a pointer to the SQL database.
type PostgresDBHandler struct {
	*sql.DB
}

// getCredentials it is a function that returns the credentials from the fiel .env to connect to the database.
func getCredentials() (string, string, string, string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	dbName := os.Getenv("DB_NAME")
	rolName := os.Getenv("ROL_NAME")
	rolPassword := os.Getenv("ROL_PASSWORD")

	return host, port, dbName, rolName, rolPassword
}

// Connect_DB it is a function that connects to the database
func Connect_DB() (*PostgresDBHandler, error) {

	host, port, dbName, rolName, rolPassword := getCredentials()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, rolName, rolPassword, dbName)
	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)

		return nil, err
	}
	fmt.Println("Successful connection to the database:", dbConn)
	DbHandler := &PostgresDBHandler{dbConn}
	var result int
	err = dbConn.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The request to the database is active")

	return DbHandler, nil
}
