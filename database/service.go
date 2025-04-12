package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"notes-app/models"
)

// Service provides methods for interacting with the database.
type Service struct {
	db *gorm.DB
}

// Connect sets up a connection to the database.
//
// This method should be called before any other methods of the Service struct.
func (svc *Service) Connect(host string, port int, user string, password string, name string, sslmode string) {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, name, sslmode,
	)

	var err error
	if svc.db, err = gorm.Open(postgres.Open(connectionString), nil); err != nil {
		panic(err)
	}
	slog.Debug("Connected to DB", slog.String("conn string", connectionString), slog.String("name", svc.db.Name()))

	err = svc.db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	slog.Debug("Migrated DB")
}

// GetDB returns the underlying Gorm DB instance.
//
// This method will panic if Connect has not been called first.
func (svc *Service) GetDB() *gorm.DB {
	if svc.db == nil {
		panic("Connect to DB first ^._.^")
	}
	return svc.db
}

// ClearAllTables clears all tables in the database.
//
// This method is intended for use in testing or development environments only.
func (svc *Service) ClearAllTables() {
	dbSession := svc.db.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true})

	dbSession.Delete(&models.User{})
}
