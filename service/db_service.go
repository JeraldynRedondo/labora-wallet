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

// GetCredentials it is a function that returns the credentials from the fiel .env to connect to the database.
func GetCredentials() (string, string, string, string, string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("host ")
	port := os.Getenv("port")
	dbName := os.Getenv("dbName")
	rolName := os.Getenv("rolName")
	rolPassword := os.Getenv("rolPassword")

	return host, port, dbName, rolName, rolPassword
}

// Connect_DB it is a function that connects to the database
func Connect_DB() (*PostgresDBHandler, error) {

	host, port, dbName, rolName, rolPassword := GetCredentials()

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, rolName, rolPassword, dbName)
	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Successful connection to the database:", dbConn)
	DbHandler := &PostgresDBHandler{dbConn}
	return DbHandler, nil
}
