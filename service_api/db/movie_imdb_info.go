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

// Get movie imdb infos
// @Summary Get movie imdb infos
// @Description Get list of all movie imdb infos
// @Tags movie_imdb_info
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /movie_imdb_info [get]
func ListMovieImdbInfoHandler(g *gin.Context) {
	infos, err := listMovieImdbInfo()

	if err != nil {
		log.Error(err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	g.JSON(http.StatusOK, gin.H{"movie_imdb_infos": infos})
}

// Add movie_imdb_info
// @Summary Add movie_imdb_info
// @Description Creates movie_imdb_info in database
// @Tags movie_imdb_info
// @Accept json
// @Produce json
// @Param movie_imdb_info body db.MovieImdbInfo true "movie_imdb_info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movie_imdb_info [post]
func AddMovieImdbInfoHandler(g *gin.Context) {
	var json MovieImdbInfo

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := addMovieImdbInfo(&json)
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

	g.JSON(http.StatusOK, gin.H{"status": "movie_imdb_info is created", "movie_imdb_info": json})
}

// Query movie_imdb_info
// @Summary Query movie_imdb_info
// @Description Shows movie_imdb_info by id
// @Tags movie_imdb_info
// @Accept json
// @Produce json
// @Param id path integer true "movie_imdb_info id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movie_imdb_info/{id} [get]
func QueryMovieImdbInfoHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := queryMovieImdbInfo(id)

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

	g.JSON(http.StatusOK, gin.H{"movie_imdb_info": info})
}

// Update movie_imdb_info
// @Summary Update movie_imdb_info
// @Description Updates movie_imdb_info specified by id
// @Tags movie_imdb_info
// @Accept json
// @Produce json
// @Param user body db.MovieImdbInfo true "movie_imdb_info"
// @Param id path integer true "movie_imdb_info id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movie_imdb_info/{id} [patch]
func UpdateMovieImdbInfoHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var json MovieImdbInfo

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = updateMovieImdbInfo(id, &json)

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

// Delete movie_imdb_info
// @Summary Delete movie_imdb_info
// @Description Delete movie_imdb_info by id
// @Tags movie_imdb_info
// @Accept json
// @Produce json
// @Param id path integer true "movie_imdb_info id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movie_imdb_info/{id} [delete]
func DeleteMovieImdbInfoHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = deleteMovieImdbInfo(id)

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

	g.JSON(http.StatusOK, gin.H{"status": "movie_imdb_info is deleted"})
}
