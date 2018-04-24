package craigslist

import (
	"time"
)

type Listing struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	URL         string            `json:"url"`
	Price       int               `json:"price"`
	Images      []Image           `json:"images"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	Location    *LatLng           `json:"location,omitempty"`
	PostedAt    *time.Time        `json:"posted_at,omitempty"`
	UpdatedAt   *time.Time        `json:"updated_at,omitempty"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Image struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type SearchResults struct {
	RangeFrom  int       `json:"range_from"`
	RangeTo    int       `json:"range_to"`
	TotalCount int       `json:"total_count"`
	Listings   []Listing `json:"listings"`
}
