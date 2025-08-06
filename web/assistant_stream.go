package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type StreamRequest struct {
	Message  string `json:"message"`
	ThreadID string `json:"threadId,omitempty"`
}

type StreamChunk struct {
	Content  string `json:"content"`
	ThreadID string `json:"threadId,omitempty"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

func HandleAssistantStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// Fallback to assistant key if general key not set
		apiKey = os.Getenv("OPENAI_ASSISTANT_API_KEY")
	}

	if apiKey == "" {
		sendStreamError(w, "OpenAI API not configured")
		return
	}

	// Parse request
	var req StreamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendStreamError(w, "Invalid request")
		return
	}

	if req.Message == "" {
		sendStreamError(w, "Message is required")
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create OpenAI client
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	// Build messages array - for now, simple conversation
	// In future, could maintain conversation history with threadID
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful AI assistant for the Blue documentation website. You help users understand Blue's features, navigate documentation, and answer questions about the platform."),
		openai.UserMessage(req.Message),
	}

	// Create streaming chat completion parameters
	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModelGPT4o,
		Messages: messages,
	}

	// Create the stream
	stream := client.Chat.Completions.NewStreaming(r.Context(), params)

	// Create a simple threadID for this session if not provided
	threadID := req.ThreadID
	if threadID == "" {
		threadID = fmt.Sprintf("thread_%d_%d", time.Now().Unix(), time.Now().Nanosecond())
	}

	// Send initial chunk with threadID
	initialChunk := StreamChunk{
		Content:  "",
		ThreadID: threadID,
		Done:     false,
	}
	sendSSEData(w, initialChunk)

	// Process the stream using the official API pattern
	for stream.Next() {
		chunk := stream.Current()
		
		// Extract content from the choice
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			streamChunk := StreamChunk{
				Content:  chunk.Choices[0].Delta.Content,
				ThreadID: threadID,
				Done:     false,
			}
			sendSSEData(w, streamChunk)
		}

		// Check if stream is finished
		if len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason != "" {
			finalChunk := StreamChunk{
				Content:  "",
				ThreadID: threadID,
				Done:     true,
			}
			sendSSEData(w, finalChunk)
			fmt.Fprintf(w, "data: [DONE]\n\n")
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			break
		}
	}

	// Check for any errors after streaming
	if err := stream.Err(); err != nil {
		errorChunk := StreamChunk{
			Error: fmt.Sprintf("Stream error: %v", err),
			Done:  true,
		}
		sendSSEData(w, errorChunk)
		fmt.Fprintf(w, "data: [DONE]\n\n")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func sendSSEData(w http.ResponseWriter, chunk StreamChunk) {
	data, err := json.Marshal(chunk)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "data: %s\n\n", string(data))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func sendStreamError(w http.ResponseWriter, message string) {
	// For streaming endpoint, we still need to send as SSE if headers are set
	if w.Header().Get("Content-Type") == "text/event-stream" {
		errorChunk := StreamChunk{
			Error: message,
			Done:  true,
		}
		sendSSEData(w, errorChunk)
		fmt.Fprintf(w, "data: [DONE]\n\n")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	} else {
		// Regular error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": message,
		})
	}
}