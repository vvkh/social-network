package config

import "os"

const (
	HostHerokuEnv = "HOST"
	PortHerokuEnv = "PORT"
)

func AdaptHerokuEnv() error {
	if os.Getenv(ServerAddressEnv) == "" {
		host := os.Getenv(HostHerokuEnv)
		port := os.Getenv(PortHerokuEnv)
		return os.Setenv(ServerAddressEnv, host+":"+port)
	}
	return nil
}
