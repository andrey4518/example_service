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

type Tag struct {
	ID      uint   `gorm:"primaryKey" json:"id" xml:"id" swaggerignore:"true"`
	UserID  uint   `form:"user_id" json:"user_id" xml:"user_id" binding:"required"`
	User    User   `gorm:"foreignKey:UserID" json:"-" swaggerignore:"true" binding:"-"`
	MovieID uint   `form:"movie_id" json:"movie_id" xml:"movie_id" binding:"required"`
	Movie   Movie  `gorm:"foreignKey:MovieID" json:"-" swaggerignore:"true" binding:"-"`
	TagText string `form:"tag_text" json:"tag_text" xml:"tag_text"  binding:"required"`
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

func addTags(tags []Tag) error {
	db, err := get_db()

	if err != nil {
		return &InternalError{Message: fmt.Sprintf("can't open database connection: %s", err.Error())}
	}

	result := db.Create(tags)
	if result.Error != nil {
		return &InternalError{Message: fmt.Sprintf("can't perform insert operation: %s", result.Error.Error())}
	}

	t := ""

	for _, m := range tags {
		t = t + strconv.Itoa(int(m.ID)) + ";"
	}

	log.Info("Insert Tags with ids: <" + t + ">")

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

// Get tags
// @Summary Get tags
// @Description Get list of all tags
// @Tags tags
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /tags [get]
func ListTagsHandler(g *gin.Context) {
	tags, err := listTags()

	if err != nil {
		log.Error(err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	g.JSON(http.StatusOK, gin.H{"tags": tags})
}

// Add tag
// @Summary Add tag
// @Description Creates tag in database
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body db.Tag true "tag info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /tags [post]
func AddTagHandler(g *gin.Context) {
	var json Tag

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := addTag(&json)
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
	g.JSON(http.StatusOK, gin.H{"status": "tag is created", "tag": json})
}

// Add tags
// @Summary Add tags
// @Description Creates tags in database
// @Tags tags
// @Accept json
// @Produce json
// @Param tags body []db.Tag true "tags info"
// @Success 200
// Failure 400
// @Failure 500
// @Router /tags/insert_batch [post]
func AddTagsHandler(g *gin.Context) {
	var json []Tag

	if err := g.ShouldBindJSON(&json); err != nil {
		body, _ := g.GetRawData()
		log.WithFields(log.Fields{"request_body": body}).Error(err)
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := addTags(json)
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

	g.JSON(http.StatusOK, gin.H{"status": "success", "tags": json})
}

// Query tag
// @Summary Query tag
// @Description Shows tag by id
// @Tags tags
// @Accept json
// @Produce json
// @Param id path integer true "tag id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /tags/{id} [get]
func QueryTagHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag, err := queryTag(id)

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
	g.JSON(http.StatusOK, gin.H{"tag": tag})
}

// Update tag
// @Summary Update tag
// @Description Updates tag info specified by id
// @Tags tags
// @Accept json
// @Produce json
// @Param user body db.Tag true "tag info"
// @Param id path integer true "tag id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /tags/{id} [patch]
func UpdateTagHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var json Tag

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = updateTag(id, &json)

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

// Delete tag
// @Summary Delete tag
// @Description Delete tag by id
// @Tags tags
// @Accept json
// @Produce json
// @Param id path integer true "tag id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /tags/{id} [delete]
func DeleteTagHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = deleteTag(id)

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

	g.JSON(http.StatusOK, gin.H{"status": "tag is deleted"})
}
