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

	// Perform a search
	result, err := craigslist.Search("chicago", opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, listing := range result.Listings {
		log.Println(listing.JSON())
	}

	// Fetch listing details
	listing, err := craigslist.GetListing("listing url")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(listing)
}
