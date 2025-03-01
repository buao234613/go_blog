package models

import (
	"WhiteBlog/common"
	"WhiteBlog/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"log"
	"sort"
	"strconv"
	"time"
)

type Article struct {
	BaseModel
	Title   string
	Content string
	//  外键
	ClassID int
}

type FormatArticle struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Class       string    `json:"class"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

var index = common.ArticleIndex

// GetArticles 获取所有文章标题
func GetArticles() ([]FormatArticle, error) {
	var articles []Article
	var formatArticles []FormatArticle
	err := config.GetDatabase().Select([]string{"ID", "title", "ClassID", "CreatedDate", "UpdatedDate"}).Find(&articles).Error
	for _, article := range articles {
		class := Class{}
		class.ID = article.ClassID
		err := class.GetClass()
		if err != nil {
			continue
		}
		formatArticle := FormatArticle{
			ID:          article.ID,
			Title:       article.Title,
			Class:       class.Name,
			CreatedDate: article.CreatedDate,
			UpdatedDate: article.UpdatedDate,
		}
		formatArticles = append(formatArticles, formatArticle)
	}
	for i, j := 0, len(formatArticles)-1; i < j; i, j = i+1, j-1 {
		formatArticles[i], formatArticles[j] = formatArticles[j], formatArticles[i]
	}
	return formatArticles, err
}

// GetArticle 获取一篇文章
func (article *Article) GetArticle() error {
	return config.GetDatabase().First(&article).Error
}

// CreateArticle 新增文章
func (article *Article) CreateArticle() error {
	return config.GetDatabase().Create(&article).Error
}

// UpdateArticle 更新文章
func (article *Article) UpdateArticle() error {
	return config.GetDatabase().Updates(&article).Error
}

// DeleteArticle 删除文章
func (article *Article) DeleteArticle() error {
	return config.GetDatabase().Delete(&article).Error
}

// GetArticlesByClass 按分类获取文章
func GetArticlesByClass(id int) ([]FormatArticle, error) {
	var formatArticles []FormatArticle
	var classes []Class
	class := Class{}
	class.ID = id
	subclasses, err := class.GetSubclasses()
	if err != nil {
		classes = append(classes, class)
	} else {
		classes = append(subclasses, class)
	}
	for _, c := range classes {
		var articles []Article
		err := config.GetDatabase().Where("class_id = ?", c.ID).Find(&articles).Error
		if err != nil {
			return formatArticles, err
		}
		for _, article := range articles {
			class := Class{}
			class.ID = article.ClassID
			err := class.GetClass()
			if err != nil {
				continue
			}
			formatArticle := FormatArticle{
				ID:          article.ID,
				Title:       article.Title,
				Class:       class.Name,
				CreatedDate: article.CreatedDate,
				UpdatedDate: article.UpdatedDate,
			}
			formatArticles = append(formatArticles, formatArticle)
		}
	}
	sort.Slice(formatArticles, func(i, j int) bool {
		return formatArticles[i].CreatedDate.After(formatArticles[j].CreatedDate)
	})
	return formatArticles, nil
}

func GetArticleByIds(ids []string) ([]FormatArticle, error) {
	var formatArticles []FormatArticle
	var articles []Article
	err := config.GetDatabase().Where("id in (?)", ids).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	for _, article := range articles {
		class := Class{}
		class.ID = article.ClassID
		err := class.GetClass()
		if err != nil {
			continue
		}
		formatArticle := FormatArticle{
			ID:          article.ID,
			Title:       article.Title,
			Class:       class.Name,
			CreatedDate: article.CreatedDate,
			UpdatedDate: article.UpdatedDate,
		}
		formatArticles = append(formatArticles, formatArticle)
	}
	return formatArticles, nil
}

// GetArticleNum 获取文章数
func GetArticleNum() int {
	var articles []Article
	err := config.GetDatabase().Find(&articles).Error
	if err != nil {
		return 0
	}
	return len(articles)
}

// IndexArticles 索引文章到es
func IndexArticles() error {

	formatArticles, err := GetArticles()
	if err != nil {
		log.Printf("IndexArticles err: %v\n", err)
		return err
	}

	esClient := config.GetClient()
	ctx := context.Background()

	exists, err := esClient.IndexExists(index).Do(ctx)
	if err != nil {
		return errors.Wrap(err, "检查索引存在失败")
	}

	if !exists {
		// 创建索引配置
		createIndex, err := esClient.CreateIndex(index).BodyString(`{
            "settings": {
                "number_of_shards": 3,
                "number_of_replicas": 2
            },
            "mappings": {
                "properties": {
                    "title": { "type": "text" },
                    "class": { "type": "text" },
                    "created_date": { "type": "date" },
                    "updated_date": { "type": "date" }
                }
            }
        }`).Do(ctx)
		if err != nil {
			return errors.Wrap(err, "索引创建失败")
		}

		if !createIndex.Acknowledged {
			return fmt.Errorf("索引创建失败，未被确认")
		}

		log.Printf("索引 %s 已创建\n", index)
	} else {
		log.Printf("索引 %s 已存在，跳过创建\n", index)
	}

	batchSize := 1000
	bulkRequest := esClient.Bulk()
	count := 0 // 用于统计当前批量请求中的文档数量

	for _, article := range formatArticles {
		req := elastic.NewBulkIndexRequest().
			Index(index).
			Id(fmt.Sprintf("%d", article.ID)). // TODO: 确保 article.ID 为数字
			Doc(article)
		bulkRequest.Add(req)
		count++

		// 每 1000 条执行一次批量请求
		if count >= batchSize {
			bulkResponse, err := bulkRequest.Do(ctx)
			if err != nil {
				log.Printf("批量请求失败: %v\n", err)
				return errors.Wrap(err, "批量请求失败")
			}

			if bulkResponse.Errors {
				log.Printf("批量请求中存在错误文档\n")
				for _, item := range bulkResponse.Items {
					for _, bulkItem := range item {
						if bulkItem.Error != nil {
							log.Printf("\t错误: %v\n", bulkItem.Error)
						}
					}
				}
				return errors.New("部分文档在批量索引时失败，请检查日志")
			}

			bulkRequest = esClient.Bulk()
			count = 0 // 重置计数器
		}
	}

	// 处理剩余文档
	if count > 0 {
		bulkResponse, err := bulkRequest.Do(ctx)
		if err != nil {
			log.Printf("批量请求失败: %v\n", err)
			return errors.Wrap(err, "批量请求失败")
		}

		if bulkResponse.Errors {
			log.Printf("批量请求中存在错误文档\n")
			for _, item := range bulkResponse.Items {
				for _, bulkItem := range item {
					if bulkItem.Error != nil {
						log.Printf("\t错误: %v\n", bulkItem.Error)
					}
				}
			}
			return errors.New("部分文档在批量索引时失败，请检查日志")
		}
	}

	log.Printf("所有文章已成功索引到 Elasticsearch\n")
	return nil
}

func SearchArticlesFromES(q string) (result []FormatArticle, err error) {
	client := config.GetClient()
	query := elastic.NewMatchQuery("title", q)
	searchedResult, err := client.Search().Index(index).Query(query).Pretty(true).Do(context.Background())
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, hit := range searchedResult.Hits.Hits {
		var article Article
		if err := json.Unmarshal(hit.Source, &article); err != nil {
			log.Printf("[SearchArticlesFromES] err : %v \n", err)
		}
		ids = append(ids, strconv.Itoa(article.ID))
	}
	log.Println("ids", ids)
	result, err = GetArticleByIds(ids)
	if err != nil {
		return nil, err
	}
	return
}
