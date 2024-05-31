package api

import (
	"log"

	"github.com/gin-gonic/gin"

	"lxdexplorer-api/config"
	"lxdexplorer-api/database"

	"github.com/gin-contrib/cors"
)

// StartAPI starts the API server
func StartAPI() {
	r := gin.Default()
	r.Use(cors.Default())

	// Load the configuration
	conf, _ := config.LoadConfig()

	// Add the routes
	r.GET("/api/v1/_healthz", health)
	r.GET("/api/v1/containers", getContainers)

	// Run the API server
	r.Run(conf.Server.Bind + ":" + conf.Server.Port)
}

func health(c *gin.Context) {

	err := database.Ping()

	if err != nil {
		log.Printf("Error pinging database: %v\n", err)

		c.JSON(500, gin.H{
			"health":   "DOWN",
			"database": "DOWN",
		})
		return
	}

	c.JSON(200, gin.H{
		"health":   "UP",
		"database": "UP",
	})
}

func getContainers(c *gin.Context) {

	containers, err := database.FindAll("containers")

	if err != nil {
		c.JSON(500, gin.H{
			"type":    "error",
			"status":  "500",
			"message": "Error fetching containers from database. Please try again later.",
		})
		return
	}

	if len(containers) == 0 {
		c.JSON(200, []gin.H{})
		return
	}

	// Remove the object_id from the containers list
	for i := range containers {
		delete(containers[i], "_id")
	}

	// Respond with the modified containers list as JSON
	c.JSON(200, containers)
}
