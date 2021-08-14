package config

import (
	"os"
	"strconv"
)

const (
	DBUrlEnv         = "DB_URL"
	AuthSecretEnv    = "AUTH_SECRET"
	ServerAddressEnv = "SERVER_ADDRESS"
	TemplatesDirEnv  = "TEMPLATES_DIR"
	BcrypCostEnv     = "BCRYPT_COST"

	BcryptDefaultCost = 10
)

type Config struct {
	DBUrl         string
	AuthSecret    string
	ServerAddress string
	TemplatesDir  string
	BcryptCost    int
}

func NewFromEnv() Config {
	bcrypCostRaw := os.Getenv(BcrypCostEnv)
	bcrypCost, err := strconv.Atoi(bcrypCostRaw)
	if err != nil {
		bcrypCost = BcryptDefaultCost
	}

	return Config{
		DBUrl:         os.Getenv(DBUrlEnv),
		AuthSecret:    os.Getenv(AuthSecretEnv),
		ServerAddress: os.Getenv(ServerAddressEnv),
		TemplatesDir:  os.Getenv(TemplatesDirEnv),
		BcryptCost:    bcrypCost,
	}
}
