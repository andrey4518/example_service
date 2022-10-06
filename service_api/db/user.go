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

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id" xml:"id" swaggerignore:"true"`
	Username string `form:"username" json:"username" xml:"username"  binding:"required"`
	Name     string `form:"name" json:"name" xml:"name"  binding:"required"`
	Sex      string `form:"sex" json:"sex" xml:"sex"  binding:"required"`
	Address  string `form:"address" json:"address" xml:"address"  binding:"required"`
	EMail    string `form:"email" json:"email" xml:"email"  binding:"required"`
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

// Get users
// @Summary Get Users
// @Description Get list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /users [get]
func ListUsersHandler(g *gin.Context) {
	users, err := listUsers()
	if err != nil {
		log.Error(err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	g.JSON(http.StatusOK, gin.H{"users": users})
}

// Add user
// @Summary Add user
// @Description Creates user in database
// @Tags users
// @Accept json
// @Produce json
// @Param user body db.User true "user info"
// @Success 200
// Failure 400
// @Failure 500
// @Router /users [post]
func AddUserHandler(g *gin.Context) {
	var json User

	if err := g.ShouldBindJSON(&json); err != nil {
		body, _ := g.GetRawData()
		log.WithFields(log.Fields{"request_body": body}).Error(err)
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := addUser(&json)
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

	g.JSON(http.StatusOK, gin.H{"status": "success", "user": json})
}

// Query user
// @Summary Query user
// @Description Shows user by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path integer true "user id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /users/{id} [get]
func QueryUserHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := queryUser(id)

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

	g.JSON(http.StatusOK, gin.H{"user": user})
}

// Update user
// @Summary Update user
// @Description Updates user info specified by id
// @Tags users
// @Accept json
// @Produce json
// @Param user body db.User true "user info"
// @Param id path integer true "user id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /users/{id} [patch]
func UpdateUserHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var json User

	if err := g.ShouldBindJSON(&json); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = updateUser(id, &json)

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

	g.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Delete user
// @Summary Delete user
// @Description Delete user by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path integer true "user id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /users/{id} [delete]
func DeleteUserHandler(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = deleteUser(id)

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

	g.JSON(http.StatusOK, gin.H{"status": "user is deleted"})
}
