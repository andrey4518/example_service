package db

import (
	"net/http"

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

func AddApiRoutes(g *gin.RouterGroup) {
	g.POST("/db/init_db", InitHandler)
	//users
	g.GET("/users", ListUsersHandler)
	g.GET("/users/:id", QueryUserHandler)
	g.POST("/users", AddUserHandler)
	g.POST("/users/insert_batch", AddUsersHandler)
	g.PATCH("/users/:id", UpdateUserHandler)
	g.DELETE("/users/:id", DeleteUserHandler)
	//movies
	g.GET("/movies", ListMoviesHandler)
	g.GET("/movies/:id", QueryMovieHandler)
	g.POST("/movies", AddMovieHandler)
	g.POST("/movies/insert_batch", AddMoviesHandler)
	g.PATCH("/movies/:id", UpdateMovieHandler)
	g.DELETE("/movies/:id", DeleteMovieHandler)
	//ratings
	g.GET("/ratings", ListRatingsHandler)
	g.GET("/ratings/:id", QueryRatingHandler)
	g.POST("/ratings", AddRatingHandler)
	g.POST("/ratings/insert_batch", AddRatingsHandler)
	g.PATCH("/ratings/:id", UpdateRatingHandler)
	g.DELETE("/ratings/:id", DeleteRatingHandler)
	//tags
	g.GET("/tags", ListTagsHandler)
	g.GET("/tags/:id", QueryTagHandler)
	g.POST("/tags", AddTagHandler)
	g.POST("/tags/insert_batch", AddTagsHandler)
	g.PATCH("/tags/:id", UpdateTagHandler)
	g.DELETE("/tags/:id", DeleteTagHandler)
	//movie imdb info
	g.GET("/movie_imdb_info", ListMovieImdbInfoHandler)
	g.GET("/movie_imdb_info/:id", QueryMovieImdbInfoHandler)
	g.POST("/movie_imdb_info", AddMovieImdbInfoHandler)
	g.POST("/movie_imdb_info/insert_batch", AddMovieImdbInfosHandler)
	g.PATCH("/movie_imdb_info/:id", UpdateMovieImdbInfoHandler)
	g.DELETE("/movie_imdb_info/:id", DeleteMovieImdbInfoHandler)
	//movie tmdb info
	g.GET("/movie_tmdb_info", ListMovieTmdbInfoHandler)
	g.GET("/movie_tmdb_info/:id", QueryMovieTmdbInfoHandler)
	g.POST("/movie_tmdb_info", AddMovieTmdbInfoHandler)
	g.POST("/movie_tmdb_info/insert_batch", AddMovieTmdbInfosHandler)
	g.PATCH("/movie_tmdb_info/:id", UpdateMovieTmdbInfoHandler)
	g.DELETE("/movie_tmdb_info/:id", DeleteMovieTmdbInfoHandler)
}
