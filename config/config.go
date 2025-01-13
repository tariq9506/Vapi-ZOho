package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq" //import postgres driver
)

// DB db object
var DB *sql.DB

// DBConfig represents db configuration
type DBConfig struct {
	Host       string
	Port       int64
	User       string
	DBName     string
	DBName2    string
	Password   string
	SSLEnabled bool
}

// BuildDBConfig builds db config object from environment variables
func BuildDBConfig() DBConfig {

	port, _ := strconv.ParseInt(os.Getenv("DBPORT"), 10, 0)
	sslEnabled, _ := strconv.ParseBool(os.Getenv("SSLENABLED"))

	dbConfig := DBConfig{
		Host:       os.Getenv("DBHOST"),
		Port:       port,
		User:       os.Getenv("DBUSER"),
		DBName:     os.Getenv("DBNAME"),
		Password:   os.Getenv("DBPASS"),
		SSLEnabled: sslEnabled,
	}
	return dbConfig
}

// DbURL get db connection string
func (dbConfig DBConfig) DbURL() string {
	log.Println("SSLEnabled", dbConfig.SSLEnabled)
	if dbConfig.SSLEnabled {
		return fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=verify-ca sslrootcert=./certs/dev/rootCA.crt",
			dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
}

func openDB() {
	var err error
	DB, err = sql.Open("postgres", BuildDBConfig().DbURL())
	if err != nil {
		log.Println("db connect error")
		log.Println("openDB: failed:", err)
		return
	}

	err = DB.Ping()
	if err != nil {
		log.Println("openDB: failed:", err)
		return
	}
}

// GetDB get db object
func GetDB() *sql.DB {
	if DB == nil {
		openDB()
	}
	return DB
}

// GetDB2 same as getDB but returns an error when it cannot connect.
func GetDB2() (*sql.DB, error) {
	DB, err := sql.Open("postgres", BuildDBConfig().DbURL())
	if err != nil {
		log.Println("GetDB2: Failed to connect to DB: " + err.Error())
		return DB, err
	}

	return DB, nil
}
