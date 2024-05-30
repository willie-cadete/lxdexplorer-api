package main

import (
	"lxdexplorer-api/api"
)

var version string

func main() {
	// print the version
	println("LXD Explorer API Version: " + version)

	// Add API routes
	api.StartAPI()
}
