package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Fprintln(os.Stderr, "OPENAI_API_KEY is not set")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 4096), 1024*1024)

	fmt.Println("OpenAI Console Chat")
	fmt.Println("Commands: /clear, /exit")

	startOpenAI(scanner)

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Input error: %v\n", err)
	}
}

func startOpenAI(scanner *bufio.Scanner) {
	const model = "gpt-5.2"
	var previousResponceID string
	client := openai.NewClient()
	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		message := strings.TrimSpace(scanner.Text())

		switch message {
		case "":
			continue
		case "/exit":
			fmt.Println("Goodbye!")
			return
		case "/clear":
			previousResponceID = ""
			fmt.Println("Conversation context cleared.")
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		params := responses.ResponseNewParams{
			Model: model,
			Store: openai.Bool(true),
			Instructions: openai.String(
				"You are a helpful assistant. " +
					"Respond in the user's language."),
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(message),
			},
		}

		if previousResponceID != "" {
			params.PreviousResponseID = openai.String(previousResponceID)
		}
		fmt.Println("AI is thinking...")
		resp, err := client.Responses.New(ctx, params)
		cancel()

		if err != nil {
			fmt.Fprintf(os.Stderr, "API error: %v\n", err)
			continue
		}

		answer := resp.OutputText()
		if answer == "" {
			fmt.Println("AI: No text response returned.")
			continue
		}

		previousResponceID = resp.ID
		fmt.Printf("\nAI: %s\n", answer)
	}
}
