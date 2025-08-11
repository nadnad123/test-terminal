package main

import (
    "bytes"
    "encoding/json"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "io"
    "net/http"
    "os"
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type CompletionRequest struct {
    Messages []Message `json:"messages"`
}

type Choice struct {
    Message struct {
        Content string `json:"content"`
    } `json:"message"`
}

type CompletionResponse struct {
    Choices []Choice `json:"choices"`
}

func main() {
    r := gin.Default()

    // Configure CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"POST"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        AllowCredentials: true,
    }))

    r.POST("/api/completion", func(c *gin.Context) {
        var req CompletionRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Prepare request to Hugging Face API
        client := &http.Client{}
        hfReqBody, _ := json.Marshal(req)
        hfReq, err := http.NewRequest("POST", "https://router.huggingface.co/v1/chat/completions", bytes.NewBuffer(hfReqBody))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HF request"})
            return
        }

        hfReq.Header.Set("Authorization", "Bearer "+os.Getenv("HF_TOKEN"))
        hfReq.Header.Set("Content-Type", "application/json")

        // Send request to Hugging Face API
        resp, err := client.Do(hfReq)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to contact HF API"})
            return
        }
        defer resp.Body.Close()

        // Read and parse response
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read HF response"})
            return
        }

        var hfResp CompletionResponse
        if err := json.Unmarshal(body, &hfResp); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse HF response"})
            return
        }

        if len(hfResp.Choices) == 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "No choices in HF response"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "response": hfResp.Choices[0].Message.Content,
        })
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    r.Run(":" + port)
}