package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/option"
)

type Req struct {
	Body string `json:"body"`
}
type Resp struct {
	Response string `json:"response"`
}

func main() {
	e := echo.New()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	// Register chat endpoint
	e.POST("/chat", chat)
	// Use CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.POST},
	}))

	log.Fatal(e.Start(os.Getenv("PORT")))
}

func chat(c echo.Context) error {
	var req Req
	if err := c.Bind(&req); err != nil {
		return err
	}
	ctx := context.Background()

	// Access API key from environment variable
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return c.JSON(http.StatusInternalServerError, "Missing API key in environment.")
	}

	// Create a new GenAI client with API key
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, "Error creating GenAI client.")
	}
	defer client.Close()

	// Use the Gemini 1.5 model
	model := client.GenerativeModel("gemini-1.5-flash")

	// Generate response based on the message
	resp, err := model.GenerateContent(ctx, genai.Text(req.Body))
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, "Error generating response.")
	}

	formattedContent := formatResponse(resp)
	return c.JSON(http.StatusOK, &Resp{
		Response: formattedContent,
	})
}

// format resposne
func formatResponse(resp *genai.GenerateContentResponse) string {
	var formattedContent strings.Builder
	if resp != nil && resp.Candidates != nil {
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					formattedContent.WriteString(fmt.Sprintf("%v", part))
				}
			}
		}
	}

	return formattedContent.String()
}
