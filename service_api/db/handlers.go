package db

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Initialize database
// @Summary Initialize database
// @Description initializes database
// @Tags db
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /db/init_db [post]
func InitHandler(g *gin.Context) {
	err := Init()
	if err != nil {
		log.Error(err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	g.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Initialize database
// @Summary Initialize test data
// @Description Initialize test data
// @Tags db
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /db/init_db_data [post]
func InitTestDataHandler(g *gin.Context) {
	err := InitTestData()
	if err != nil {
		log.Error(err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	g.JSON(http.StatusOK, gin.H{"status": "success"})
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

	g.JSON(http.StatusOK, gin.H{"status": "movie is created"})
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

func AddApiRoutes(g *gin.RouterGroup) {
	g.POST("/db/init_db", InitHandler)
	g.POST("/db/init_db_data", InitTestDataHandler)
	//users
	g.GET("/users", ListUsersHandler)
	g.GET("/users/:id", QueryUserHandler)
	g.POST("/users", AddUserHandler)
	g.PATCH("/users/:id", UpdateUserHandler)
	g.DELETE("/users/:id", DeleteUserHandler)
	//movies
	g.GET("/movies", ListMoviesHandler)
	g.GET("/movies/:id", QueryMovieHandler)
	g.POST("/movies", AddMovieHandler)
	g.PATCH("/movies/:id", UpdateMovieHandler)
	g.DELETE("/movies/:id", DeleteMovieHandler)
}
