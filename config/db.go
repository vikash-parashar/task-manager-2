package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// my local postgres db credentials, please do change if needed
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "task-manager-two"
)

var DB *sql.DB

// @Summary Initialize the database connection
// @Description Establishes a connection to the PostgreSQL database
// @ID init-db
// @Produce json
// @Success 200 {string} string "Successfully connected to the database"
// @Failure 500 {object} string "Internal server error"
// @Router /init [get]
func InitDB() error {
	var err error
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")
	return nil
}
