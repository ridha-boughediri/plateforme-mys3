package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}
}
