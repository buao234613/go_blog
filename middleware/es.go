package middleware

// Elasticsearch
import (
	"WhiteBlog/models"
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

type ES struct {
	esClient *elastic.Client
}

var es *ES
var once sync.Once

func GetESClient() *elastic.Client {
	if es == nil {
		once.Do(func() {
			es = &ES{}
		})
	}
	return es.esClient
}
func (e *ES) InitElasticSearch(url string) error {
	if e.esClient != nil {
		return nil // 已经初始化过了
	}
	esClient, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return fmt.Errorf("elasticsearch connection failed: %w", err)
	}
	e.esClient = esClient
	return nil
}

// 数据库中数据存储到es
func (e *ES) IndexDataFromDB(db *gorm.DB) error {
	if e.esClient == nil {
		return fmt.Errorf("elasticsearch client is nil")
	}
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	articles, err := models.GetArticles()
	if err != nil {
		return err
	}
	for _, article := range articles {
		_, err := e.esClient.Index().
			Index("articles").
			Id(strconv.Itoa(article.ID)).
			BodyJson(article).
			Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}
