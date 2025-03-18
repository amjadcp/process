package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type EnvStruct struct {
	GROQ_API_URL string
	GROQ_API_KEY string
	GROQ_MODEL   string
}

var (
	once sync.Once
	Env  *EnvStruct
)

func init() {
	// Initialize configuration once
	once.Do(func() {
		// Load .env file
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Warning: Error loading .env file (using system env): %v", err)
		}

		Env = &EnvStruct{
			GROQ_API_URL: os.Getenv("GROQ_API_URL"),
			GROQ_API_KEY: os.Getenv("GROQ_API_KEY"),
			GROQ_MODEL:   os.Getenv("GROQ_MODEL"),
		}
	})
}