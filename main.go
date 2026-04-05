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
	"github.com/manifoldco/promptui"
	"github.com/ollama/ollama/api"
	"github.com/schollz/progressbar/v3"
)

func main() {
	ollamaServer := flag.String("url", "http://127.0.0.1:11434", "URL of the Ollama server")
	flag.Parse()

	parsedURL, err := url.Parse(*ollamaServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid URL %s\n", *ollamaServer)
		os.Exit(1)
	}
	client := api.NewClient(parsedURL, http.DefaultClient)

	model, err := selectModel(client)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error selecting model:", err)
		os.Exit(1)
	}

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

		response, err := query(modelContext, model, client)
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

func query(modelContext []api.Message, model string, client *api.Client) (string, error) {
	req := api.ChatRequest{
		Model:    model,
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

// selectModel let user select model from the list of models available on the Ollama server. It returns the name of the selected model.
func selectModel(client *api.Client) (string, error) {
	ctx := context.Background()

	resp, err := client.List(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list models: %w", err)
	}

	if len(resp.Models) == 0 {
		return "", fmt.Errorf("no models available on the server")
	}

	names := make([]string, len(resp.Models))
	for i, m := range resp.Models {
		names[i] = m.Name
	}

	prompt := promptui.Select{
		Label: "Select model",
		Items: names,
	}

	_, selected, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return selected, nil
}
