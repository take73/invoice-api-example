package rdb

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"

	// MySQL migrate driver
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const defaultDBName = "upsider"

func NewDB() (*gorm.DB, error) {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = defaultDBName
	}

	return connectDB(dbName, 5)
}

func connectDB(dbName string, retryCount int) (*gorm.DB, error) {
	conn, err := gorm.Open(gormmysql.Open(loadConfig(dbName).FormatDSN()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		if retryCount == 0 {
			return nil, fmt.Errorf(": %w", err)
		}
		log.Printf("db connect retry %d\n", retryCount)
		time.Sleep(2 * time.Second)
		return connectDB(dbName, retryCount-1)
	}
	return conn, nil
}

func loadConfig(dbName string) *mysql.Config {
	return &mysql.Config{
		DBName:               dbName,
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Addr:                 os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Net:                  os.Getenv("DB_NET"),
		Loc:                  time.UTC,
		ParseTime:            true,
		AllowNativePasswords: true,
	}
}
