package db

import (
	"fmt"
	"os"
	"pose-service/models"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var (
		connectionName = mustGetenv("CLOUDSQL_CONNECTION_NAME")
		user           = mustGetenv("CLOUDSQL_USER")
		dbName         = os.Getenv("CLOUDSQL_DATABASE_NAME")
		password       = os.Getenv("CLOUDSQL_PASSWORD")
		socket         = os.Getenv("CLOUDSQL_SOCKET_PREFIX")
	)

	// /cloudsql is used on App Engine.
	if socket == "" {
		socket = "/cloudsql"
	}

	connString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", connectionName, user, dbName, password)
	database, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsqlpostgres",
		DSN:        connString,
	}))
	if err != nil {
		panic(err)
	}

	db = database

	err = db.AutoMigrate(&models.Pose{})
	if err != nil {
		return
	}
}

func GetDB() *gorm.DB {
	return db
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(k + " environment variable not set.")
	}
	return v
}
