package db

import (
	"errors"
	notifier "example/service/api/notifier"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Rating struct {
	ID      uint    `gorm:"primaryKey" json:"id" xml:"id" swaggerignore:"true"`
	UserID  uint    `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
	User    User    `gorm:"foreignKey:UserID" json:"-" swaggerignore:"true" binding:"-"`
	MovieID uint    `form:"movie_id" json:"movie_id" xml:"movie_id" binding:"required"`
	Movie   Movie   `gorm:"foreignKey:MovieID" json:"-" swaggerignore:"true" binding:"-"`
	Rating  float32 `form:"rating" json:"rating" xml:"rating" binding:"required"`
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

// Get ratings
// @Summary Get ratings
// @Description Get list of all ratings
// @Tags ratings
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /ratings [get]
func ListRatingsHandler(g *gin.Context) {
	ratings, err := listRatings()

	if err != nil {
		log.Error(err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	g.JSON(http.StatusOK, gin.H{"ratings": ratings})
}

// Add rating
// @Summary Add rating
// @Description Creates rating in database
// @Tags ratings
// @Accept json
// @Produce json
// @Param rating body db.Rating true "rating info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /ratings [post]
func AddRatingHandler(g *gin.Context) {
	var json Rating

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := addRating(&json)
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

	g.JSON(http.StatusOK, gin.H{"status": "rating is created", "rating": json})
}

// Query rating
// @Summary Query rating
// @Description Shows rating by id
// @Tags ratings
// @Accept json
// @Produce json
// @Param id path integer true "rating id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /ratings/{id} [get]
func QueryRatingHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rating, err := queryRating(id)

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

	g.JSON(http.StatusOK, gin.H{"rating": rating})
}

// Update rating
// @Summary Update rating
// @Description Updates rating info specified by id
// @Tags ratings
// @Accept json
// @Produce json
// @Param user body db.Rating true "rating info"
// @Param id path integer true "rating id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /ratings/{id} [patch]
func UpdateRatingHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var json Rating

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = updateRating(id, &json)

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

// Delete rating
// @Summary Delete rating
// @Description Delete rating by id
// @Tags ratings
// @Accept json
// @Produce json
// @Param id path integer true "rating id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /ratings/{id} [delete]
func DeleteRatingHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = deleteRating(id)

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

	g.JSON(http.StatusOK, gin.H{"status": "rating is deleted"})
}
