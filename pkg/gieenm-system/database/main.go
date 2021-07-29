package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"

	// postgreql driver
	_ "github.com/lib/pq"
)

var db *sqlx.DB

//Init ...
func Init() {
	dbInfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
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

//GetDB ...
func GetDB() *sqlx.DB {
	return db
}

//RedisClient ...
var RedisClient *redis.Client

//InitRedis ...
func InitRedis(params ...string) {
	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	db, _ := strconv.Atoi(params[0])

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       db,
	})
}

//GetRedis ...
func GetRedis() *redis.Client {
	return RedisClient
}
