package db

import (
	"fmt"
	"strconv"

	"example/service/api/config"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _db *gorm.DB

func get_db() (*gorm.DB, error) {
	if _db == nil {
		db, err := gorm.Open(postgres.Open(config.GetDbConnectionString()), &gorm.Config{})
		if err != nil {
			return nil, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
		}
		_db = db
	}
	return _db, nil
}

func Init() error {
	db, err := get_db()
	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	db.AutoMigrate(
		&Movie{},
		&User{},
		&Rating{},
		&Tag{},
		&MovieImdbInfo{},
		&MovieTmdbInfo{},
	)
	log.Info("Database initialized")
	return nil
}

func InitTestData() error {
	for i := 0; i < 3; i++ {
		err := addUser(&User{Name: "user" + strconv.Itoa(i)})
		if err != nil {
			return fmt.Errorf("inserting user error: %w", err)
		}
	}

	for i := 0; i < 3; i++ {
		err := addMovie(&Movie{Name: "movie" + strconv.Itoa(i)})
		if err != nil {
			return fmt.Errorf("inserting user error: %w", err)
		}
	}

	log.Info("Test data inserted")

	return nil
}
