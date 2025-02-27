package article

import (
	"WhiteBlog/common"
	"WhiteBlog/models"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	query := c.Query("query")
	result, err := models.SearchArticlesFromES(query)
	if err != nil {
		common.NotFound(c, "Not found matched result!")
		return
	}
	common.Render(c, "searchpage.html", gin.H{"result": result})
}
