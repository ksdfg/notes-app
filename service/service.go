package service

import (
	"gorm.io/gorm"
	"notes-app/database"
)

// DBOpts is a struct that contains common options for CRUD methods.
//
// It can be used to pass a separate DB instance to run CRUD operations in transactions that are controlled externally.
type DBOpts struct {
	db *gorm.DB
}

// Service defines the common fields and methods for all services.
type Service struct {
	DBService database.Service
}

// getDB returns a DB instance for CRUD operations, either from the provided options or the default DB instance.
func (svc Service) getDB(opts *DBOpts) *gorm.DB {
	if opts != nil && opts.db != nil {
		return opts.db
	}
	return svc.DBService.GetDB()
}
