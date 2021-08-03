package postgres

import (
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

var models = []interface{}{
	&account{},
	&application{},
	&device{},
	&tag{},
}

func CreateRepository(db *gorm.DB) (*repository, error) {
	repo := &repository{
		db: db,
	}

	rawDB, err := db.DB()
	if err != nil {
		return repo, errors.Wrap(err, "failed to get raw sql from gorm")
	}

	// Enable UUID Extensions
	_, err = rawDB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public`)
	if err != nil {
		return repo, errors.Wrap(err, "failed to migrate uuid extension")
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		return repo, errors.Wrap(err, "failed to auto migrate models")
	}
	return repo, nil
}

func getProcessedDBError(err error) error  {
	if err == gorm.ErrRecordNotFound {
		return errors.WithKindCtx(err, "", errors.NotFound, nil)
	}
	return err
}
