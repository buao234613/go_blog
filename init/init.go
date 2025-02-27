package init

import (
	"WhiteBlog/config"
	"WhiteBlog/models"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func init() {
	jsonConfig, err := os.Open("config.json")
	if err != nil {
		return
	}
	defer jsonConfig.Close()
	err = json.NewDecoder(jsonConfig).Decode(&config.TheConfig)
	//  数据库
	databaseConfig := config.GetConfig().DatabaseConfig
	switch config.GetConfig().DatabaseConfig.Driver {
	case "mysql":
		//初始化MySQL
		dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true",
			databaseConfig.Username,
			databaseConfig.Password,
			databaseConfig.Host,
			databaseConfig.Port,
			databaseConfig.Database,
			databaseConfig.Charset)
		log.Println(dsn)
		conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}
		config.Database = conn
	}
	//  迁移表
	err = config.Database.AutoMigrate(
		&models.Class{},
		&models.Article{},
		&models.Image{},
		&models.Link{},
		&models.Poet{},
		&models.File{},
	)
	// 初始化 es
	esConfig := config.GetConfig().ESConfig
	url := fmt.Sprintf("http://%s:%s", esConfig.Host, esConfig.Port)
	esClient, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		log.Fatalf("elasticsearch 客户端初始化失败: %v", err)
	}
	config.Client = esClient
	err = models.IndexArticles()
	if err != nil {
		log.Fatalf("IndexArticles error: %v", err)
	}
}
