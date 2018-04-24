package craigslist

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyResults(t *testing.T) {
	results, err := ParseSearchResults(strings.NewReader(""))

	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 0, results.RangeFrom)
	assert.Equal(t, 0, results.RangeTo)
	assert.Equal(t, 0, results.TotalCount)
	assert.Empty(t, results.Listings)
}

func TestParseResults(t *testing.T) {
	file, err := os.Open("./fixtures/results.html")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	results, err := ParseSearchResults(file)
	assert.NotNil(t, results)
	assert.Equal(t, 1, results.RangeFrom)
	assert.Equal(t, 120, results.RangeTo)
	assert.Equal(t, 3000, results.TotalCount)
	assert.Equal(t, 120, len(results.Listings))

	listing := results.Listings[0]
	assert.Equal(t, "6569794207", listing.Id)
	assert.Equal(t, "2011 Porsche Cayenne Turbo", listing.Title)
	assert.Equal(t, "https://chicago.craigslist.org/wcl/cto/d/2011-porsche-cayenne-turbo/6569794207.html", listing.URL)
	assert.Empty(t, listing.Description)
	assert.Equal(t, 39000, listing.Price)
	assert.Equal(t, 10, len(listing.Images))
	assert.NotNil(t, listing.PostedAt)
	assert.NotNil(t, listing.UpdatedAt)
}
