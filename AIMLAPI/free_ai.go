package AIMLAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const AIMLAPIBaseURL = "https://api.aimlapi.com/v1"
const AIMLAPIKey = "3008ba967ecb488bb2cfef65536ea2dd" // Replace with your actual API key

// RequestPayload represents the structure of the API request.
type RequestPayload struct {
	Model       string        `json:"model"`
	Messages    []MessageItem `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// MessageItem represents a message in the conversation.
type MessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ResponsePayload represents the structure of the API response.
type ResponsePayload struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func callAIMLAPI() {
	// Set the prompt details
	systemPrompt := "You are a helpful assistant."
	userPrompt := "Write a short description of GoLang."

	// Prepare the request payload
	requestPayload := RequestPayload{
		Model: "mistralai/Mistral-7B-Instruct-v0.2",
		Messages: []MessageItem{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
		MaxTokens:   256,
	}

	// Serialize the payload to JSON
	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		log.Fatalf("Error serializing request payload: %v", err)
	}

	// Create an HTTP POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", AIMLAPIBaseURL), bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+AIMLAPIKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API Error: %s", string(responseData))
	}

	var responsePayload ResponsePayload
	if err := json.Unmarshal(responseData, &responsePayload); err != nil {
		log.Fatalf("Error parsing JSON response: %v", err)
	}

	// Extract and display the AI response
	if len(responsePayload.Choices) > 0 {
		answer := responsePayload.Choices[0].Message.Content
		log.Println("User:", userPrompt)
		log.Println("AI:", answer)
	} else {
		log.Println("No response from AI.")
	}
}
