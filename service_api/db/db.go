package db

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Movie struct {
	ID   uint   `gorm:"primaryKey" swaggerignore:"true"`
	Name string `form:"name" json:"name" xml:"name"  binding:"required"`
}

type User struct {
	ID   uint   `gorm:"primaryKey" swaggerignore:"true"`
	Name string `form:"name" json:"name" xml:"name"  binding:"required"`
}

var _db *gorm.DB

func get_db() (*gorm.DB, error) {
	if _db == nil {
		db, err := gorm.Open(postgres.Open(viper.GetString("common.pg_connection_string")), &gorm.Config{})
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

	db.AutoMigrate(&Movie{}, &User{})
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

func listUsers() ([]User, error) {
	var users []User
	db, err := get_db()

	if err != nil {
		return users, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Find(&users)

	if result.Error != nil {
		return users, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", err.Error())}
	}

	return users, nil
}

func addUser(u *User) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Create(u)
	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform insert operation: %s", result.Error.Error())}
	}

	log.Info("Insert User with id: <" + strconv.Itoa(int(u.ID)) + ">")

	return nil
}

func queryUser(id int) (User, error) {
	db, err := get_db()

	var user User

	if err != nil {
		return user, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Where("id = ?", id).First(&user)

	if result.Error != nil {
		return user, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return user, &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return user, nil
}

func updateUser(id int, user *User) error {
	db, err := get_db()
	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	var data User
	result := db.Where("id = ?", id).First(&data)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	result = db.Model(&data).Select("*").Omit("id").Updates(user)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform update operation: %s", result.Error.Error())}
	}

	return nil
}

func deleteUser(id int) error {
	db, err := get_db()
	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Delete(&User{}, id)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform delete operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return nil
}

func listMovies() ([]Movie, error) {
	db, err := get_db()

	var movies []Movie

	if err != nil {
		return movies, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Find(&movies)

	if result.Error != nil {
		return movies, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", err.Error())}
	}

	return movies, nil
}

func addMovie(m *Movie) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Create(m)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform insert operation: %s", result.Error.Error())}
	}

	log.Info("Insert Movie with id: <" + strconv.Itoa(int(m.ID)) + ">")

	return nil
}

func queryMovie(id int) (Movie, error) {
	db, err := get_db()
	var movie Movie

	if err != nil {
		return movie, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Where("id = ?", id).First(&movie)

	if result.Error != nil {
		return movie, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return movie, &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return movie, nil
}

func updateMovie(id int, movie *Movie) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	var data Movie
	result := db.Where("id = ?", id).First(&data)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	result = db.Model(&data).Select("*").Omit("id").Updates(movie)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform update operation: %s", result.Error.Error())}
	}

	return nil
}

func deleteMovie(id int) error {
	db, err := get_db()
	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Delete(&Movie{}, id)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform delete operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return nil
}
