package ai

import (
	"fmt"

	"github.com/amjadcp/process/config"
)

type Ollama struct {
	URL string
	Model string
	APIKEY string
}


func (o *Ollama)Chat()  {
	o.URL = "https://api.groq.com/openai/v1/chat/completions"
	fmt.Println(config.Env.GROQ_API_KEY)
}