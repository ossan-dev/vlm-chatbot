package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ollama/ollama/api"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 1. get Ollama client
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(err)
	}

	if err := client.Heartbeat(ctx); err != nil {
		panic(err)
	}

	// 2. preparing the image
	imgBytes, err := os.ReadFile("/home/ossan/Projects/vlm-chatbot/booking_resized.png")
	if err != nil {
		panic(err)
	}

	// 3. adding tiles to the context
	messages := []api.Message{
		{
			Role: "system",
			Content: "You are a literal OCR engine. Your only job is to extract text from images. " +
				"Rule 1: Never guess or invent data. " +
				"Rule 2: If a field is present but unreadable, write [UNREADABLE]. " +
				"Rule 3: If a field is missing, write [NOT FOUND]. " +
				"Rule 4: Output only the extracted text without any conversational preamble.",
		},
		{Role: "user",
			Content: "Extract the text from this image exactly as it appears",
			Images:  []api.ImageData{imgBytes},
		},
	}

	// set LLM options
	options := map[string]any{
		"temperature": 0.0,  // force the most likely token
		"top_p":       0.1,  // only consider high-probability tokens
		"num_ctx":     4096, // ensure enough context for the image tokens
	}

	falsePtr := false
	req := &api.ChatRequest{
		Model:    "qwen3-vl:235b-cloud", // or qwen3-vl:235b-cloud to use the cloud-based one
		Messages: messages,
		Stream:   &falsePtr,
		Options:  options,
	}

	// 4. setup Ollama response
	chatFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		if resp.Done {
			req.Messages = append(req.Messages, resp.Message)
		}
		fmt.Printf(" - replied in: %.2fs\n", resp.Metrics.TotalDuration.Seconds())
		return nil
	}

	// kick-off conversation
	err = client.Chat(ctx, req, chatFunc)
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	prompts := []string{
		"Extract the total cost of the reservation.",
		"Extract the check-in and the check-out dates.",
		"Extract the address of the accommodation.",
	}

	for _, p := range prompts {
		fmt.Printf("\nUser > %q", p)
		req.Messages = append(req.Messages, api.Message{
			Role:    "user",
			Content: p,
		})
		fmt.Printf("\nAssistant> ")
		err = client.Chat(ctx, req, chatFunc)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}
	}
}
