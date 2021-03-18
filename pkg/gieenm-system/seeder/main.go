package seeder

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// InitPostgres ...
func InitPostgres() {
	dbInfo := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	var err error

	if db, err = ConnectDB(dbInfo); err != nil {
		log.Fatal(err)
	}
}

//ConnectDB ...
func ConnectDB(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}

// Init ...
func Init() {
	InitPostgres()
}

// Seeder ...
func Seeder() error {
	FirewallSeeder()
	AnnouncementSeeder()
	StudentCandySeeder()

	return nil
}
