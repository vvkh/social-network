package config

import "os"

const (
	DBUrlEnv         = "DB_URL"
	AuthSecretEnv    = "AUTH_SECRET"
	ServerAddressEnv = "SERVER_ADDRESS"
	TemplatesDirEnv  = "TEMPLATES_DIR"
)

type Config struct {
	DBUrl         string
	AuthSecret    string
	ServerAddress string
	TemplatesDir  string
}

func NewFromEnv() Config {
	return Config{
		DBUrl:         os.Getenv(DBUrlEnv),
		AuthSecret:    os.Getenv(AuthSecretEnv),
		ServerAddress: os.Getenv(ServerAddressEnv),
		TemplatesDir:  os.Getenv(TemplatesDirEnv),
	}
}
