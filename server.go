package craigslist

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func searchHandler(c *gin.Context) {
	site := c.Query("site")
	if site == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "site is required"})
		return
	}

	opts := SearchOptions{}
	if err := c.BindQuery(&opts); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	ids := strings.Split(site, ",")
	if len(ids) > 1 {
		result, err := MultiSearch(ids, opts)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, result)
		return
	}

	result, err := Search(site, opts)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func StartServer(addr string) error {
	router := gin.Default()
	router.GET("/search", searchHandler)
	return router.Run(addr)
}
