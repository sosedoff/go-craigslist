package craigslist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchOptionsQuery(t *testing.T) {
	q := SearchOptions{Query: "search query"}.query()
	assert.Equal(t, 1, len(q))
	assert.Equal(t, "search query", q.Get("query"))

	q = SearchOptions{
		Query:            "query",
		TitlesOnly:       true,
		HasImage:         true,
		PostedToday:      true,
		BundleDuplicates: true,
		IncludeNearby:    true,
		MinPrice:         10,
		MaxPrice:         20,
		Skip:             100,
		Params:           map[string]string{"foo": "bar"},
	}.query()

	assert.Equal(t, 10, len(q))
	assert.Equal(t, "query", q.Get("query"))
	assert.Equal(t, "T", q.Get("srchType"))
	assert.Equal(t, "1", q.Get("hasPic"))
	assert.Equal(t, "1", q.Get("postedToday"))
	assert.Equal(t, "1", q.Get("bundleDuplicates"))
	assert.Equal(t, "1", q.Get("searchNearby"))
	assert.Equal(t, "10", q.Get("min_price"))
	assert.Equal(t, "20", q.Get("max_price"))
	assert.Equal(t, "100", q.Get("s"))
}
