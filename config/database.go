package config

import (
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Charset  string `json:"charset"`
	//Local    string
}

var Database *gorm.DB

func GetDatabase() *gorm.DB {
	return Database
}
