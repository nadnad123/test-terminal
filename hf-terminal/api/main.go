package handler

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "os"
)

// Message represents a single message in the Hugging Face API request
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// CompletionRequest represents the incoming JSON payload
type CompletionRequest struct {
    Messages []Message `json:"messages"`
}

// Choice represents a choice in the Hugging Face API response
type Choice struct {
    Message struct {
        Content string `json:"content"`
    } `json:"message"`
}

// CompletionResponse represents the Hugging Face API response
type CompletionResponse struct {
    Choices []Choice `json:"choices"`
}

// Handler is the Vercel serverless function entry point
func Handler(w http.ResponseWriter, r *http.Request) {
    // Ensure the request is a POST
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse the incoming JSON payload
    var req CompletionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Prepare request to Hugging Face API
    client := &http.Client{}
    hfReqBody, _ := json.Marshal(req)
    hfReq, err := http.NewRequest("POST", "https://router.huggingface.co/v1/chat/completions", bytes.NewBuffer(hfReqBody))
    if err != nil {
        http.Error(w, "Failed to create HF request", http.StatusInternalServerError)
        return
    }

    // Set headers for Hugging Face API
    hfReq.Header.Set("Authorization", "Bearer "+os.Getenv("HF_TOKEN"))
    hfReq.Header.Set("Content-Type", "application/json")

    // Send request to Hugging Face API
    resp, err := client.Do(hfReq)
    if err != nil {
        http.Error(w, "Failed to contact HF API", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read and parse response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Failed to read HF response", http.StatusInternalServerError)
        return
    }

    var hfResp CompletionResponse
    if err := json.Unmarshal(body, &hfResp); err != nil {
        http.Error(w, "Failed to parse HF response", http.StatusInternalServerError)
        return
    }

    if len(hfResp.Choices) == 0 {
        http.Error(w, "No choices in HF response", http.StatusInternalServerError)
        return
    }

    // Prepare response for the client
    response := map[string]string{
        "response": hfResp.Choices[0].Message.Content,
    }

    // Set CORS headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Content-Type", "application/json")

    // Send JSON response
    json.NewEncoder(w).Encode(response)
}
