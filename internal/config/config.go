package config

import "os"

//var localdbConfig = "postgresql://postgres:1234@localhost:5432/effective?sslmode=disable"

func GetDBConfig() string {
	if s := os.Getenv("PG_DSN"); s != "" {
		return s
	}
	localdbConfig := os.Getenv("LOCAL_DB_CONFIG")

	return localdbConfig
}
