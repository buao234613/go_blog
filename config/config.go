package config

type Config struct {
	AuthConfig     `json:"authConfig"`
	DatabaseConfig `json:"databaseConfig"`
	ESConfig       `json:"esConfig"`
	User           `json:"user"`
	Site           `json:"site"`
}

var TheConfig Config

func GetConfig() Config {
	return TheConfig
}
