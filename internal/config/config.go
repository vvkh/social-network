package config

import "os"

type Config struct {
	DBUrl         string
	AuthSecret    string
	ServerAddress string
	TemplatesDir  string
}

func NewFromEnv() Config {
	return Config{
		DBUrl:         os.Getenv("DB_URL"),
		AuthSecret:    os.Getenv("AUTH_SECRET"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		TemplatesDir:  os.Getenv("TEMPLATES_DIR"),
	}
}
