package main

import "lxdexplorer-api/fetcher"

func main() {

	// Prepare TTL Indexes
	fetcher.AddLXDTTLs()

	for {
		fetcher.Run()
	}

}
