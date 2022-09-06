package db

import (
	"fmt"
	"strconv"

	"example/service/api/config"

	pq "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Movie struct {
	ID      uint           `gorm:"primaryKey" json:"id" xml:"id" swaggerignore:"true"`
	Name    string         `form:"name" json:"name" xml:"name" binding:"required"`
	Imdb_Id uint           `form:"imdb_id" json:"imdb_id" xml:"imdb_id" binding:"required"`
	Tmdb_Id uint           `form:"tmdb_id" json:"tmdb_id" xml:"tmdb_id" binding:"required"`
	Genres  pq.StringArray `gorm:"type:varchar(64)[]" form:"genres" json:"genres" xml:"genres" binding:"required" swaggertype:"array,string"`
}

type MovieImdbInfo struct {
	ID            uint           `gorm:"primaryKey" json:"id" xml:"id" swaggerignore:"true"`
	MovieId       uint           `form:"movie_id" json:"movie_id" xml:"movie_id"  binding:"required"`
	Genres        pq.StringArray `gorm:"type:text[]" form:"genres" json:"genres" xml:"genres" binding:"required" swaggertype:"array,string"`
	OriginalTitle string         `form:"original_title" json:"original_title" xml:"original_title"`
	Runtimes      pq.StringArray `gorm:"type:text[]" form:"runtimes" json:"runtimes" xml:"runtimes" binding:"required" swaggertype:"array,string"`
	Countries     pq.StringArray `gorm:"type:text[]" form:"countries" json:"countries" xml:"countries" binding:"required" swaggertype:"array,string"`
	Rating        float32        `form:"rating" json:"rating" xml:"rating" binding:"required"`
	Votes         uint           `form:"votes" json:"votes" xml:"votes"  binding:"required"`
	PlotOutline   string         `gorm:"type:text" form:"plot_outline" json:"plot_outline" xml:"plot_outline"`
	Languages     pq.StringArray `gorm:"type:text[]" form:"languages" json:"languages" xml:"languages" binding:"required" swaggertype:"array,string"`
	Year          uint           `form:"year" json:"year" xml:"year"  binding:"required"`
	Kind          string         `form:"kind" json:"kind" xml:"kind"`
	Plot          pq.StringArray `gorm:"type:text[]" form:"plot" json:"plot" xml:"plot" binding:"required" swaggertype:"array,string"`
	Synopsis      pq.StringArray `gorm:"type:text[]" form:"synopsis" json:"synopsis" xml:"synopsis" binding:"required" swaggertype:"array,string"`
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id" xml:"id" swaggerignore:"true"`
	Username string `form:"username" json:"username" xml:"username"  binding:"required"`
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Sex      string `form:"sex" json:"sex" xml:"sex"  binding:"required"`
	Address  string `form:"address" json:"address" xml:"address"  binding:"required"`
	EMail    string `form:"email" json:"email" xml:"email"  binding:"required"`
}

type Rating struct {
	ID      uint    `gorm:"primaryKey" json:"id" xml:"id" swaggerignore:"true"`
	UserID  uint    `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
	User    User    `gorm:"foreignKey:UserID" json:"-" swaggerignore:"true" binding:"-"`
	MovieID uint    `form:"movie_id" json:"movie_id" xml:"movie_id" binding:"required"`
	Movie   Movie   `gorm:"foreignKey:MovieID" json:"-" swaggerignore:"true" binding:"-"`
	Rating  float32 `form:"rating" json:"rating" xml:"rating" binding:"required"`
}

type Tag struct {
	ID      uint   `gorm:"primaryKey" json:"id" xml:"id" swaggerignore:"true"`
	UserID  uint   `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
	User    User   `gorm:"foreignKey:UserID" json:"-" swaggerignore:"true" binding:"-"`
	MovieID uint   `form:"movie_id" json:"movie_id" xml:"movie_id" binding:"required"`
	Movie   Movie  `gorm:"foreignKey:MovieID" json:"-" swaggerignore:"true" binding:"-"`
	TagText string `form:"tag_text" json:"tag_text" xml:"tag_text"  binding:"required"`
}

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

func listRatings() ([]Rating, error) {
	db, err := get_db()

	var ratings []Rating

	if err != nil {
		return ratings, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Find(&ratings)

	if result.Error != nil {
		return ratings, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", err.Error())}
	}

	return ratings, nil
}

func addRating(r *Rating) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Create(r)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform insert operation: %s", result.Error.Error())}
	}

	log.Info("Insert Rating with id: <" + strconv.Itoa(int(r.ID)) + ">")

	return nil
}

func queryRating(id int) (Rating, error) {
	db, err := get_db()
	var rating Rating

	if err != nil {
		return rating, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Where("id = ?", id).First(&rating)

	if result.Error != nil {
		return rating, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return rating, &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return rating, nil
}

func updateRating(id int, rating *Rating) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	var data Rating
	result := db.Where("id = ?", id).First(&data)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	result = db.Model(&data).Select("*").Omit("id").Updates(rating)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform update operation: %s", result.Error.Error())}
	}

	return nil
}

func deleteRating(id int) error {
	db, err := get_db()
	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Delete(&Rating{}, id)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform delete operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return nil
}

func listMovieImdbInfo() ([]MovieImdbInfo, error) {
	db, err := get_db()

	var infos []MovieImdbInfo

	if err != nil {
		return infos, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Find(&infos)

	if result.Error != nil {
		return infos, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", err.Error())}
	}

	return infos, nil
}

func addMovieImdbInfo(i *MovieImdbInfo) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Create(i)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform insert operation: %s", result.Error.Error())}
	}

	log.Info("Insert MovieImdbInfo with id: <" + strconv.Itoa(int(i.ID)) + ">")

	return nil
}

func queryMovieImdbInfo(id int) (MovieImdbInfo, error) {
	db, err := get_db()
	var info MovieImdbInfo

	if err != nil {
		return info, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Where("id = ?", id).First(&info)

	if result.Error != nil {
		return info, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return info, &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return info, nil
}

func updateMovieImdbInfo(id int, info *MovieImdbInfo) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	var data MovieImdbInfo
	result := db.Where("id = ?", id).First(&data)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	result = db.Model(&data).Select("*").Omit("id").Updates(info)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform update operation: %s", result.Error.Error())}
	}

	return nil
}

func deleteMovieImdbInfo(id int) error {
	db, err := get_db()
	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Delete(&MovieImdbInfo{}, id)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform delete operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return nil
}

func listTags() ([]Tag, error) {
	db, err := get_db()

	var tags []Tag

	if err != nil {
		return tags, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Find(&tags)

	if result.Error != nil {
		return tags, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", err.Error())}
	}

	return tags, nil
}

func addTag(t *Tag) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Create(t)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform insert operation: %s", result.Error.Error())}
	}

	log.Info("Insert Tag with id: <" + strconv.Itoa(int(t.ID)) + ">")

	return nil
}

func queryTag(id int) (Tag, error) {
	db, err := get_db()
	var tag Tag

	if err != nil {
		return tag, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Where("id = ?", id).First(&tag)

	if result.Error != nil {
		return tag, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return tag, &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return tag, nil
}

func updateTag(id int, tag *Tag) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	var data Tag
	result := db.Where("id = ?", id).First(&data)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	result = db.Model(&data).Select("*").Omit("id").Updates(tag)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform update operation: %s", result.Error.Error())}
	}

	return nil
}

func deleteTag(id int) error {
	db, err := get_db()
	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Delete(&Tag{}, id)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform delete operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return nil
}
