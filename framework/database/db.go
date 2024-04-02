package database

import (
	"log"

	"github.com/brandaoplaster/encoder/domain"
	"github.com/jinzhu/gorm"
)

type Database struct {
	Db          *gorm.DB
	Dsn         string
	DsnTest     string
	DbType      string
	DbTypeTest  string
	Debug       bool
	AutoMigrate bool
	Env         string
}

func NewDatabase() *Database {
	return &Database{}
}

func NewDatabaseTest() *gorm.DB {
	dbInstance := NewDatabase()
	dbInstance.Env = "test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory:"
	dbInstance.Debug = true
	dbInstance.AutoMigrate = true

	connection, error := dbInstance.Connect()

	if error != nil {
		log.Fatalf("Error database test => %v", error)
	}

	return connection
}

func (database *Database) Connect() (*gorm.DB, error) {
	var erro error

	if database.Env == "test" {
		database.Db, erro = gorm.Open(database.DbTypeTest, database.DsnTest)
	} else {
		database.Db, erro = gorm.Open(database.DbType, database.Dsn)
	}

	if erro != nil {
		return nil, erro
	}

	if database.Debug {
		database.Db.LogMode(true)
	}

	if database.AutoMigrate {
		database.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
	}

	return database.Db, nil
}
