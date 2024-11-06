package server

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func InitDatabase(config *viper.Viper) *sql.DB {
	connString := config.GetString("database.connection_string")
	maxIdleConns := config.GetInt("database.max_idle_connections")
	maxOpenconns := config.GetInt("database.max_open_connections")
	connMaxLifetime := config.GetDuration("database.connection_max_lifetime")
	driverName := config.GetString("database.driver_name")
	if connString == "" {
		log.Fatalf("Database connection string is missing.")
	}
	dbHandler, err := sql.Open(driverName, connString)
	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}
	dbHandler.SetMaxIdleConns(maxIdleConns)
	dbHandler.SetMaxOpenConns(maxOpenconns)
	dbHandler.SetConnMaxLifetime(connMaxLifetime)
	err = dbHandler.Ping()
	if err != nil {
		dbHandler.Close()
		log.Fatalf("Error while validating database: %v", err)
	}
	return dbHandler
}
