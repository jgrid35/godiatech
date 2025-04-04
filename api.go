package main

import (
	"github.com/gin-gonic/gin"
)

func api() {
	router := gin.Default()
	router.GET("/movies", getMovies)

	router.Run("localhost:8080")
}

func getMovies(c *gin.Context) {
	db := connectToDatabase()
	movies, err := getAllMovies(db)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch movies" + err.Error()})
		return
	}
	c.JSON(200, movies)
}
