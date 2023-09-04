package database

import (
	"fmt"
	"log"
	"os"

	"github.com/melihcanclk/docker-postgres-go-rest-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	Db *gorm.DB
}

var DB DBInstance

func DropUnusedColumns(dst interface{}) {

	stmt := &gorm.Statement{DB: DB.Db}
	stmt.Parse(dst)
	fields := stmt.Schema.Fields
	columns, _ := DB.Db.Debug().Migrator().ColumnTypes(dst)

	for i := range columns {
		found := false
		for j := range fields {
			if columns[i].Name() == fields[j].DBName {
				found = true
				break
			}
		}
		if !found {
			DB.Db.Migrator().DropColumn(dst, columns[i].Name())
		}
	}
}

func ConnectDB() {
	dataSourceName := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	Db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database, %v", err)
		os.Exit(2)
	}

	log.Println("Connected to database")
	Db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations...")

	// get models that we'll migrate, if we wanna use multiple migration
	// we can define var models = []interface{}{&User{}, &Product{}, &Order{}}
	// then db.Automigrate(models...)
	// Db.AutoMigrate(&models.Fact{})

	Db.AutoMigrate(&models.Question{}, &models.User{}, &models.Answer{})

	DB = DBInstance{
		Db,
	}

}
