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

type MovieTmdbInfo struct {
	ID            uint           `gorm:"primaryKey" json:"id" xml:"id" swaggerignore:"true"`
	MovieId       uint           `form:"movie_id" json:"movie_id" xml:"movie_id"  binding:"required"`
	Adult         *bool          `form:"adult" json:"adult" xml:"adult"  binding:"required"`
	Genres        pq.StringArray `gorm:"type:text[]" form:"genres" json:"genres" xml:"genres" binding:"required" swaggertype:"array,string"`
	HomePage      string         `form:"homepage" json:"homepage" xml:"homepage"`
	OriginalTitle string         `form:"original_title" json:"original_title" xml:"original_title"`
	Overview      string         `form:"overview" json:"overview" xml:"overview"`
	Popularity    float32        `form:"popularity" json:"popularity" xml:"popularity" binding:"required"`
	Runtime       uint           `form:"runtime" json:"runtime" xml:"runtime" binding:"required"`
	Tagline       string         `form:"tagline" json:"tagline" xml:"tagline"`
	Title         string         `form:"title" json:"title" xml:"title"`
	VoteAverage   float32        `form:"vote_average" json:"vote_average" xml:"vote_average" binding:"required"`
	VoteCount     uint           `form:"vote_count" json:"vote_count" xml:"vote_count" binding:"required"`
	Keywords      pq.StringArray `gorm:"type:text[]" form:"keywords" json:"keywords" xml:"keywords" binding:"required" swaggertype:"array,string"`
	VideoURLs     pq.StringArray `gorm:"type:text[]" form:"video_urls" json:"video_urls" xml:"video_urls" binding:"required" swaggertype:"array,string"`
}

func listMovieTmdbInfo() ([]MovieTmdbInfo, error) {
	db, err := get_db()

	var infos []MovieTmdbInfo

	if err != nil {
		return infos, &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Find(&infos)

	if result.Error != nil {
		return infos, &InternalError{Message: fmt.Sprintf("can't perform query operation: %s", err.Error())}
	}

	return infos, nil
}

func addMovieTmdbInfo(i *MovieTmdbInfo) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Create(i)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform insert operation: %s", result.Error.Error())}
	}

	log.Info("Insert MovieTmdbInfo with id: <" + strconv.Itoa(int(i.ID)) + ">")

	return nil
}

func addMovieTmdbInfos(infos []MovieTmdbInfo) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Create(infos)
	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform insert operation: %s", result.Error.Error())}
	}

	t := ""

	for _, m := range infos {
		t = t + strconv.Itoa(int(m.ID)) + ";"
	}

	log.Info("Insert MovieTmdbInfos with ids: <" + t + ">")

	return nil
}

func queryMovieTmdbInfo(id int) (MovieTmdbInfo, error) {
	db, err := get_db()
	var info MovieTmdbInfo

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

func updateMovieTmdbInfo(id int, info *MovieTmdbInfo) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	var data MovieTmdbInfo

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

func deleteMovieTmdbInfo(id int) error {
	db, err := get_db()
	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Delete(&MovieTmdbInfo{}, id)

	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform delete operation: %s", result.Error.Error())}
	}

	if result.RowsAffected == 0 {
		return &QueryConditionError{Message: fmt.Sprintf("can't find object by this id <%d>", id)}
	}

	return nil
}

// Get movie tmdb infos
// @Summary Get movie tmdb infos
// @Description Get list of all movie tmdb infos
// @Tags movie_tmdb_info
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /movie_tmdb_info [get]
func ListMovieTmdbInfoHandler(g *gin.Context) {
	infos, err := listMovieTmdbInfo()

	if err != nil {
		log.Error(err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	g.JSON(http.StatusOK, gin.H{"movie_tmdb_infos": infos})
}

// Add movie_tmdb_info
// @Summary Add movie_tmdb_info
// @Description Creates movie_tmdb_info in database
// @Tags movie_tmdb_info
// @Accept json
// @Produce json
// @Param movie_tmdb_info body db.MovieTmdbInfo true "movie_tmdb_info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movie_tmdb_info [post]
func AddMovieTmdbInfoHandler(g *gin.Context) {
	var json MovieTmdbInfo

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := addMovieTmdbInfo(&json)
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

	g.JSON(http.StatusOK, gin.H{"status": "movie_tmdb_info is created", "movie_tmdb_info": json})
}

// Add movie_tmdb_infos
// @Summary Add movie_tmdb_infos
// @Description Creates movie_tmdb_infos in database
// @Tags movie_tmdb_info
// @Accept json
// @Produce json
// @Param movie_tmdb_infos body []db.MovieTmdbInfo true "movie_tmdb_infos"
// @Success 200
// Failure 400
// @Failure 500
// @Router /movie_tmdb_info/insert_batch [post]
func AddMovieTmdbInfosHandler(g *gin.Context) {
	var json []MovieTmdbInfo

	if err := g.ShouldBindJSON(&json); err != nil {
		body, _ := g.GetRawData()
		log.WithFields(log.Fields{"request_body": body}).Error(err)
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := addMovieTmdbInfos(json)
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
	for _, u := range json {
		notifier.ObjectCreationNotificationChannel <- u
	}

	g.JSON(http.StatusOK, gin.H{"status": "success", "movie_tmdb_infos": json})
}

// Query movie_tmdb_info
// @Summary Query movie_tmdb_info
// @Description Shows movie_tmdb_info by id
// @Tags movie_tmdb_info
// @Accept json
// @Produce json
// @Param id path integer true "movie_tmdb_info id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movie_tmdb_info/{id} [get]
func QueryMovieTmdbInfoHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := queryMovieTmdbInfo(id)

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

	g.JSON(http.StatusOK, gin.H{"movie_tmdb_info": info})
}

// Update movie_tmdb_info
// @Summary Update movie_tmdb_info
// @Description Updates movie_tmdb_info specified by id
// @Tags movie_tmdb_info
// @Accept json
// @Produce json
// @Param user body db.MovieTmdbInfo true "movie_tmdb_info"
// @Param id path integer true "movie_tmdb_info id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movie_tmdb_info/{id} [patch]
func UpdateMovieTmdbInfoHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var json MovieTmdbInfo

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = updateMovieTmdbInfo(id, &json)

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

// Delete movie_tmdb_info
// @Summary Delete movie_tmdb_info
// @Description Delete movie_tmdb_info by id
// @Tags movie_tmdb_info
// @Accept json
// @Produce json
// @Param id path integer true "movie_tmdb_info id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /movie_tmdb_info/{id} [delete]
func DeleteMovieTmdbInfoHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = deleteMovieTmdbInfo(id)

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

	g.JSON(http.StatusOK, gin.H{"status": "movie_tmdb_info is deleted"})
}
