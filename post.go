package json_server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

var PostStore = map[int]*Post{}
var lock = sync.RWMutex{}

func GetPosts(c *gin.Context) {
	/*
		sort := c.Param("_sort")
		order := c.Param("_order")
		start := c.Param("_start")
		end := c.Param("_end")
	*/
	lock.RLock()
	posts := make([]Post, len(PostStore))
	i := 0
	for _, post := range PostStore {
		posts[i] = *post
		i++
	}
	lock.RUnlock()

	c.Header("X-Total-Count", strconv.Itoa(len(PostStore)))
	c.Header("Access-Control-Expose-Headers", "X-Total-Count")
	c.JSON(http.StatusOK, posts)
}
