package config

import (
	"github.com/olivere/elastic/v7"
)

type ESConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var Client *elastic.Client

func GetClient() *elastic.Client {
	return Client
}
