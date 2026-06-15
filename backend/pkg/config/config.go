package config

import (
	"os"

	"github.com/joho/godotenv"
)

func GetVar(name string) string {
	godotenv.Load("../../../.env")
	return os.Getenv(name)
}
