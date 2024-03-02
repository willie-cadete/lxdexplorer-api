package api

import (
	"github.com/gin-gonic/gin"

	"lxdexplorer-api/config"
)

// StartAPI starts the API server
func StartAPI() {
	r := gin.Default()

	// Load the configuration
	conf, _ := config.LoadConfig()

	// Run the API server
	r.Run(conf.Server.Bind + ":" + conf.Server.Port)
}
