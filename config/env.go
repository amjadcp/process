package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type EnvStruct struct {
	GROQ_API_URL   string
	GROQ_API_KEY   string
	GROQ_MODEL     string
	OLLAMA_API_URL string
	OLLAMA_API_KEY string
	OLLAMA_MODEL   string
	AI_SERVICE string
}

var (
	once sync.Once
	Env  *EnvStruct
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func init() {
	// Initialize configuration once
	once.Do(func() {
		// Load .env file
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Warning: Error loading .env file (using system env): %v", err)
		}

		Env = &EnvStruct{
			GROQ_API_URL: getEnv("GROQ_API_URL", "https://api.groq.com/openai/v1/chat/completions"),
			GROQ_API_KEY: getEnv("GROQ_API_KEY", ""),
			GROQ_MODEL:   getEnv("GROQ_MODEL", "llama-3.3-70b-versatile"),
			OLLAMA_API_URL: getEnv("OLLAMA_API_URL", "http://localhost:11434/api/chat"),
			OLLAMA_API_KEY: getEnv("OLLAMA_API_KEY", ""),
			OLLAMA_MODEL: getEnv("OLLAMA_MODEL", "llama3.2:1b"),
			AI_SERVICE: getEnv("AI_SERVICE", "groq"),
		}
	})
}
