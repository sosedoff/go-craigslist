package main

import (
	"log"

	"github.com/sosedoff/go-craigslist"
)

func main() {
	opts := craigslist.SearchOptions{
		Category: "cto", // cars+trucks
		Query:    "honda",
		HasImage: true,
		MinPrice: 10000,
		MaxPrice: 20000,
	}

	result, err := craigslist.Search("chicago", opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, listing := range result.Listings {
		log.Println(listing.JSON())
	}
}
