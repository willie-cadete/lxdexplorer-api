package main

import (
	"lxdexplorer-api/api"
	"lxdexplorer-api/fetcher"
)

func main() {

	// Create a new fetcher
	go fetcher.StartFetcher()

}
