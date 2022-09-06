package db

import (
	"errors"
	notifier "example/service/api/notifier"
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

	rating, err := queryMovieImdbInfo(id)

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

	g.JSON(http.StatusOK, gin.H{"movie_imdb_info": rating})
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
	//ratings
	g.GET("/ratings", ListRatingsHandler)
	g.GET("/ratings/:id", QueryRatingHandler)
	g.POST("/ratings", AddRatingHandler)
	g.PATCH("/ratings/:id", UpdateRatingHandler)
	g.DELETE("/ratings/:id", DeleteRatingHandler)
	//tags
	g.GET("/tags", ListTagsHandler)
	g.GET("/tags/:id", QueryTagHandler)
	g.POST("/tags", AddTagHandler)
	g.PATCH("/tags/:id", UpdateTagHandler)
	g.DELETE("/tags/:id", DeleteTagHandler)
	//movie imdb info
	g.GET("/movie_imdb_info", ListMovieImdbInfoHandler)
	g.GET("/movie_imdb_info/:id", QueryMovieImdbInfoHandler)
	g.POST("/movie_imdb_info", AddMovieImdbInfoHandler)
	g.PATCH("/movie_imdb_info/:id", UpdateMovieImdbInfoHandler)
	g.DELETE("/movie_imdb_info/:id", DeleteMovieImdbInfoHandler)
}
