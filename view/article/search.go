package article

import (
	"WhiteBlog/common"
	"WhiteBlog/models"
	"github.com/gin-gonic/gin"
	"log"
)

func Search(c *gin.Context) {
	query := c.Query("query")
	log.Printf("query : %v \n", query)
	result, err := models.SearchArticlesFromES(query)
	if err != nil {
		common.NotFound(c, "Not found matched result!")
		return
	}
	log.Println("result : ", result)
	common.Render(c, "searched.html", gin.H{"articles": result})
}
