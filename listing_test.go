package craigslist

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInvalidListing(t *testing.T) {
	listing, err := ParseListing(strings.NewReader(""))
	assert.Error(t, errInvalidListing, err)
	assert.Nil(t, listing)
}

func TestParseEmptyListing(t *testing.T) {
	file, err := os.Open("./fixtures/empty_listing.html")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	listing, err := ParseListing(file)
	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, "6569794207", listing.Id)
}

func TestParseListing(t *testing.T) {
	file, err := os.Open("./fixtures/listing2.html")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	listing, err := ParseListing(file)
	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, "6520962081", listing.Id)
	assert.Equal(t, "2014 PORSCHE PANAMERA", listing.Title)
	assert.Equal(t, "https://chicago.craigslist.org/sox/cto/d/2014-porsche-panamera/6520962081.html", listing.URL)
	assert.NotEmpty(t, listing.Description)
	assert.Equal(t, 49980, listing.Price)
	assert.Equal(t, 24, len(listing.Images))
	assert.NotNil(t, listing.PostedAt)
	assert.NotNil(t, listing.UpdatedAt)
	assert.Equal(t, 41.7447, listing.Location.Lat)
	assert.Equal(t, -87.7699, listing.Location.Lng)

	attrs := map[string]string{
		"cylinders":    "6 cylinders",
		"drive":        "4wd",
		"fuel":         "gas",
		"odometer":     "72574",
		"title_status": "clean",
		"transmission": "automatic",
		"vin":          "WP0BB2A74EL062226",
	}
	for k, v := range attrs {
		assert.Equal(t, v, listing.Attributes[k])
	}
}
