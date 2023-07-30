package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Get returns a pointer to the database
func Get() *gorm.DB {
	return db
}

// BeginTxn begins a database transaction and return a pointer to the database
func BeginTxn() *gorm.DB {
	return db.Begin()
}

// RecoverTxn recovers from a panicked transaction and rolls back
func RecoverTxn(txn *gorm.DB) {
	if r := recover(); r != nil {
		txn.Rollback()
	}
}

// CommitTxn commits the transaction and rolls back the database in case of any errors
func CommitTxn(txn *gorm.DB) error {
	if commitErr := txn.Commit().Error; commitErr != nil {
		txn.Rollback()
		return commitErr
	}
	return nil
}

// Connect opens a connection to the database
func Connect(url string, maxIdleConnections, maxOpenConnections int) error {
	var err error
	var gormdb *sql.DB
	sqlDB, err := sql.Open("nrpostgres", url)
	if err != nil {
		return err
	}
	db, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		DisableAutomaticPing: false,
	})
	if err != nil {
		return err
	}
	// Get database connection handle [*sql.DB](http://golang.org/pkg/database/sql/#DB)
	gormdb, err = db.DB()
	if err != nil {
		return err
	}

	// Then you could invoke `*sql.DB`'s functions with it
	err = gormdb.Ping()
	if err != nil {
		return err
	}
	//gormdb
	//db.LogMode(false)
	gormdb.SetMaxIdleConns(maxIdleConnections)
	gormdb.SetMaxOpenConns(maxOpenConnections)
	//db.SingularTable(false)
	return nil
}

// Close closes the database
func Close() {
	sqldb, _ := db.DB()
	_ = sqldb.Close()
}
