package api

import (
	"github.com/gin-gonic/gin"

	"lxdexplorer-api/config"
	"lxdexplorer-api/database"
)

// StartAPI starts the API server
func StartAPI() {
	r := gin.Default()

	// Load the configuration
	conf, _ := config.LoadConfig()

	// Add the routes
	r.GET("/api/v1/ping", ping)
	r.GET("/api/v1/containers", getContainers)

	// Run the API server
	r.Run(conf.Server.Bind + ":" + conf.Server.Port)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func getContainers(c *gin.Context) {

	containers, err := database.FindAll("containers")

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error fetching containers",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": containers,
	})
}
