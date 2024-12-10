package testutils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"

	// MySQL migrate driver
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func SetupTestDB(dbName string) *gorm.DB {
	projectDir, err := GetProjectDir()
	if err != nil {
		panic(err)
	}

	db, err := connectDB("", 5)
	if err != nil {
		panic(err)
	}
	dbConn, _ := db.DB()
	defer dbConn.Close()

	if createDatabase(db, dbName); err != nil {
		panic(err)
	}

	migrateDir := projectDir + "/db/migrations"
	migration(dbName, migrateDir)

	testDB, err := connectDB(dbName, 5)
	if err != nil {
		panic(err)
	}
	return testDB
}

func createDatabase(db *gorm.DB, dbName string) error {
	if err := db.Exec("drop database if exists " + dbName).Error; err != nil {
		return err
	}
	return db.Exec("create database " + dbName).Error
}

func connectDB(dbName string, retryCount int) (*gorm.DB, error) {
	str := loadConfig(dbName).FormatDSN()
	// conn, err := gorm.Open(gormmysql.Open(loadConfig(dbName).FormatDSN()), &gorm.Config{
	// 	NamingStrategy: schema.NamingStrategy{SingularTable: true},
	// })
	conn, err := gorm.Open(gormmysql.Open(str), &gorm.Config{
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

func ExecSQLFile(db *gorm.DB, filePath string) error {
	// Read the SQL file
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read SQL file %s: %v", filePath, err)
		return err
	}

	// Convert to string and split by ';' to handle multiple statements
	sqlStatements := strings.Split(string(sqlBytes), ";")

	// Execute each SQL statement individually
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue // skip empty statements
		}
		if err := db.Exec(stmt).Error; err != nil {
			log.Printf("Failed to execute SQL statement: %v", err)
			return err
		}
	}

	log.Printf("Successfully executed SQL file %s", filePath)
	return nil
}

func migration(dbName string, path string) {
	dbConfig := loadConfig(dbName)
	fqn := dbConfig.FormatDSN()
	m, err := migrate.New(
		"file://"+path,
		"mysql://"+fqn,
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		panic(err)
	}
	if err, dbError := m.Close(); err != nil {
		panic(err)
	} else if dbError != nil {
		panic(dbError)
	}
}
