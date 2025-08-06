package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type AssistantRequest struct {
	Message  string `json:"message"`
	ThreadID string `json:"threadId,omitempty"`
}

type AssistantResponse struct {
	Response string `json:"response"`
	ThreadID string `json:"threadId"`
	Error    string `json:"error,omitempty"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIThreadRequest struct {
	Messages []OpenAIMessage `json:"messages,omitempty"`
}

type OpenAIRunRequest struct {
	AssistantID string `json:"assistant_id"`
}

type OpenAIThreadResponse struct {
	ID string `json:"id"`
}

type OpenAIMessageResponse struct {
	ID      string `json:"id"`
	Content []struct {
		Text struct {
			Value string `json:"value"`
		} `json:"text"`
	} `json:"content"`
}

type OpenAIRunResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type OpenAIMessagesResponse struct {
	Data []OpenAIMessageResponse `json:"data"`
}

func HandleAssistant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get API key and Assistant ID from environment
	apiKey := os.Getenv("OPENAI_ASSISTANT_API_KEY")
	assistantID := os.Getenv("OPENAI_ASSISTANT_ID")

	if apiKey == "" || assistantID == "" {
		sendAssistantError(w, "AI Assistant not configured")
		return
	}

	// Parse request
	var req AssistantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendAssistantError(w, "Invalid request")
		return
	}

	if req.Message == "" {
		sendAssistantError(w, "Message is required")
		return
	}

	// Create or use existing thread
	threadID := req.ThreadID
	if threadID == "" {
		thread, err := createThread(apiKey)
		if err != nil {
			sendAssistantError(w, "Failed to create conversation thread")
			return
		}
		threadID = thread.ID
	}

	// Add message to thread
	if err := addMessage(apiKey, threadID, req.Message); err != nil {
		sendAssistantError(w, "Failed to send message")
		return
	}

	// Run the assistant
	run, err := runAssistant(apiKey, threadID, assistantID)
	if err != nil {
		sendAssistantError(w, "Failed to process message")
		return
	}

	// Wait for completion
	response, err := waitForCompletion(apiKey, threadID, run.ID)
	if err != nil {
		sendAssistantError(w, "Failed to get response")
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AssistantResponse{
		Response: response,
		ThreadID: threadID,
	})
}

func createThread(apiKey string) (*OpenAIThreadResponse, error) {
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/threads", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create thread: %d", resp.StatusCode)
	}

	var thread OpenAIThreadResponse
	if err := json.NewDecoder(resp.Body).Decode(&thread); err != nil {
		return nil, err
	}

	return &thread, nil
}

func addMessage(apiKey, threadID, message string) error {
	body := map[string]string{
		"role":    "user",
		"content": message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", 
		fmt.Sprintf("https://api.openai.com/v1/threads/%s/messages", threadID),
		bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add message: %d", resp.StatusCode)
	}

	return nil
}

func runAssistant(apiKey, threadID, assistantID string) (*OpenAIRunResponse, error) {
	body := OpenAIRunRequest{
		AssistantID: assistantID,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs", threadID),
		bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to run assistant: %d - %s", resp.StatusCode, string(body))
	}

	var run OpenAIRunResponse
	if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
		return nil, err
	}

	return &run, nil
}

func waitForCompletion(apiKey, threadID, runID string) (string, error) {
	for i := 0; i < 30; i++ { // Max 30 seconds
		time.Sleep(1 * time.Second)

		req, err := http.NewRequest("GET",
			fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs/%s", threadID, runID),
			nil)
		if err != nil {
			return "", err
		}

		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("OpenAI-Beta", "assistants=v2")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var run OpenAIRunResponse
		if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
			return "", err
		}

		if run.Status == "completed" {
			// Get the latest message
			return getLatestMessage(apiKey, threadID)
		} else if run.Status == "failed" || run.Status == "cancelled" || run.Status == "expired" {
			return "", fmt.Errorf("run failed with status: %s", run.Status)
		}
	}

	return "", fmt.Errorf("timeout waiting for response")
}

func getLatestMessage(apiKey, threadID string) (string, error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://api.openai.com/v1/threads/%s/messages?limit=1", threadID),
		nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var messages OpenAIMessagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return "", err
	}

	if len(messages.Data) > 0 && len(messages.Data[0].Content) > 0 {
		return messages.Data[0].Content[0].Text.Value, nil
	}

	return "", fmt.Errorf("no response from assistant")
}

func sendAssistantError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(AssistantResponse{
		Error: message,
	})
}