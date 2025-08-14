package movie

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (controller *Controller) GetMovieById(ctx *gin.Context) {
	id := ctx.Param("id")

	movieData, err := controller.tmdbService.GetMovieById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie data"})
		return
	}

	ctx.JSON(200, movieData)
}

func (controller *Controller) SearchForMovie(ctx *gin.Context) {
	query := ctx.Query("query")
	if query == "" {
		ctx.JSON(400, gin.H{"error": "Query parameter 'query' is required"})
		return
	}
	
	pageStr := ctx.DefaultQuery("page", "1")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid page number"})
		return
	}

	results, err := controller.tmdbService.SearchForMovie(query, page)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to search for movies"})
		return
	}

	ctx.JSON(200, results)
}