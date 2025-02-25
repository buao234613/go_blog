package article

import (
	"WhiteBlog/middleware"
	"WhiteBlog/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
	"strconv"
)

func Search(c *gin.Context) {
	k := c.Query("query")
	client := middleware.GetESClient()
	ids, err := searchFromElasticsearch(client, k)
	if err != nil {
		return
	}
	articles, err := models.GetArticleByIds(ids)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}

// searchFromElasticsearch 根据关键词从 Elasticsearch 中查询文档 ID
func searchFromElasticsearch(client *elastic.Client, keyword string) ([]string, error) {
	query := elastic.NewMatchQuery("name", keyword)
	searchResult, err := client.Search().
		Index("articles").
		Query(query).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found a total of %d Users\n", searchResult.TotalHits())

	var ids []string
	for _, hit := range searchResult.Hits.Hits {
		var article models.FormatArticle
		err := json.Unmarshal(hit.Source, &article)
		if err != nil {
			log.Printf("Unmarshal failed: %s", err)
			continue
		}
		ids = append(ids, strconv.Itoa(article.ID))
	}
	return ids, nil
}
