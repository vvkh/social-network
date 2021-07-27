package heroku

import "os"

func ConvertEnv() error {
	if os.Getenv("SERVER_ADDRESS") == "" {
		host := os.Getenv("HOST")
		port := os.Getenv("PORT")
		return os.Setenv("SERVER_ADDRESS", host+":"+port)
	}
	return nil
}
