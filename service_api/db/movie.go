package db

import (
	"errors"
	notifier "example/service/api/notifier"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pq "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Movie struct {
	ID      uint           `gorm:"primaryKey" json:"id" xml:"id" swaggerignore:"true"`
	Name    string         `form:"name" json:"name" xml:"name" binding:"required"`
	Imdb_Id uint           `form:"imdb_id" json:"imdb_id" xml:"imdb_id" binding:"required"`
	Tmdb_Id uint           `form:"tmdb_id" json:"tmdb_id" xml:"tmdb_id" binding:"required"`
	Genres  pq.StringArray `gorm:"type:varchar(64)[]" form:"genres" json:"genres" xml:"genres" binding:"required" swaggertype:"array,string"`
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

// Get movies
// @Summary Get movies
// @Description Get list of all movies
// @Tags movies
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /movies [get]
func ListMoviesHandler(g *gin.Context) {
	movies, err := listMovies()

	if err != nil {
		log.Error(err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	g.JSON(http.StatusOK, gin.H{"movies": movies})
}

// Add movie
// @Summary Add movie
// @Description Creates movie in database
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body db.Movie true "movie info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movies [post]
func AddMovieHandler(g *gin.Context) {
	var json Movie

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := addMovie(&json)
	if err != nil {
		switch {
		case errors.As(err, &intErr):
			log.Error(err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		case errors.As(err, &qCondErr):
			log.Error(err)
			g.JSON(http.StatusBadRequest, gin.H{"error": err})
		default:
			log.Error(err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	notifier.ObjectCreationNotificationChannel <- json

	g.JSON(http.StatusOK, gin.H{"status": "movie is created", "movie": json})
}

// Query movie
// @Summary Query movie
// @Description Shows movie by id
// @Tags movies
// @Accept json
// @Produce json
// @Param id path integer true "movie id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movies/{id} [get]
func QueryMovieHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie, err := queryMovie(id)

	if err != nil {
		switch {
		case errors.As(err, &intErr):
			log.Error(err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		case errors.As(err, &qCondErr):
			log.Error(err)
			g.JSON(http.StatusBadRequest, gin.H{"error": err})
		default:
			log.Error(err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	g.JSON(http.StatusOK, gin.H{"movie": movie})
}

// Update movie
// @Summary Update movie
// @Description Updates movie info specified by id
// @Tags movies
// @Accept json
// @Produce json
// @Param user body db.Movie true "movie info"
// @Param id path integer true "movie id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movies/{id} [patch]
func UpdateMovieHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var json Movie

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = updateMovie(id, &json)

	if err != nil {
		switch {
		case errors.As(err, &intErr):
			log.Error(err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		case errors.As(err, &qCondErr):
			log.Error(err)
			g.JSON(http.StatusBadRequest, gin.H{"error": err})
		default:
			log.Error(err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	g.JSON(http.StatusOK, gin.H{"status": "sucess"})
}

// Delete movie
// @Summary Delete movie
// @Description Delete movie by id
// @Tags movies
// @Accept json
// @Produce json
// @Param id path integer true "movie id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movies/{id} [delete]
func DeleteMovieHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = deleteMovie(id)

	if err != nil {
		switch {
		case errors.As(err, &intErr):
			log.Error(err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		case errors.As(err, &qCondErr):
			log.Error(err)
			g.JSON(http.StatusBadRequest, gin.H{"error": err})
		default:
			log.Error(err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	g.JSON(http.StatusOK, gin.H{"status": "movie is deleted"})
}
