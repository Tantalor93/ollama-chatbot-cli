package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ollama/ollama/api"
	"github.com/schollz/progressbar/v3"
)

func main() {
	ollamaServer := flag.String("url", "http://127.0.0.1:11434", "URL of the Ollama server")
	flag.Parse()
	
	parsedURL, err := url.Parse(*ollamaServer)
	if err != nil {
		panic(err)
	}
	client := api.NewClient(parsedURL, http.DefaultClient)

	var modelContext []api.Message = []api.Message{
		{
			Role: "system",
			Content: "You are a helpful assistant. " +
				"Answer the user's questions to the best of your ability, " +
				"but keep it concise as this is a CLI application.",
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	printUserPrompt()
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if len(input) == 0 {
			continue
		}
		fmt.Println("-----")

		modelContext = append(modelContext, api.Message{
			Role:    "user",
			Content: input,
		})

		response, err := query(modelContext, client)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		printModelPrompt()
		fmt.Println(response)
		fmt.Println("-----")

		modelContext = append(modelContext, api.Message{
			Role:    "assistant",
			Content: response,
		})
		printUserPrompt()
	}
}

func query(modelContext []api.Message, client *api.Client) (string, error) {
	req := api.ChatRequest{
		Model:    "llama3.2",
		Messages: modelContext,
		Stream:   new(bool),
	}

	progress := progressbar.NewOptions(
		-1,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionClearOnFinish(),
	)
	defer progress.Finish()

	var result string
	err := client.Chat(context.Background(), &req, func(resp api.ChatResponse) error {
		result = resp.Message.Content
		return nil
	})

	return result, err
}

func printUserPrompt() {
	color.New(color.FgGreen).Print("> ")
}

func printModelPrompt() {
	color.New(color.FgYellow).Print("< ")
}
