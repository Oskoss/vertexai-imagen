package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	auth "golang.org/x/oauth2/google"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/generate", generateImg)
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}

func generateImg(c *gin.Context) {
	// Check incomming request for validity
	var reqJson VertexAiRequest
	if err := c.ShouldBindJSON(&reqJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Setup GCP Auth
	var token *oauth2.Token
	ctx := context.Background()
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
	credentials, err := auth.FindDefaultCredentials(ctx, scopes...)
	if err == nil {
		fmt.Println("found default credentials!")
		token, err = credentials.TokenSource.Token()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	projectId := os.Getenv("PROJECT_ID")

	// Create Request
	url := fmt.Sprintf("https://us-central1-aiplatform.googleapis.com/v1/projects/%s/locations/us-central1/publishers/google/models/imagegeneration:predict", projectId)
	reqBytes, err := json.Marshal(reqJson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vertexReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Setup Headers
	token.SetAuthHeader(vertexReq)
	vertexReq.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Send the request
	resp, err := http.DefaultClient.Do(vertexReq)
	if err != nil {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": body})
		return
	}

	// Check the response status code
	if resp.StatusCode != 200 {
		c.JSON(http.StatusBadRequest, resp.Body)
		return
	}

	// Read and return the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result VertexAiResponse
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

type VertexAiRequest struct {
	Instances []struct {
		Prompt string `json:"prompt,omitempty"`
	} `json:"instances,omitempty"`
	Parameters struct {
		SampleCount       int    `json:"sampleCount,omitempty"`
		StorageURI        string `json:"storageUri,omitempty"`
		Seed              int    `json:"seed,omitempty"`
		NegativePrompt    string `json:"negativePrompt,omitempty"`
		DisablePersonFace bool   `json:"disablePersonFace,omitempty"`
		Mode              string `json:"mode,omitempty"`
		SampleImageSize   string `json:"sampleImageSize,omitempty"`
		IncludeRaiReason  bool   `json:"includeRaiReason,omitempty"`
	} `json:"parameters,omitempty"`
}

type VertexAiResponse struct {
	Predictions []struct {
		BytesBase64Encoded string `json:"bytesBase64Encoded,omitempty"`
		Image              string `json:"image,omitempty"`
		MimeType           string `json:"mimeType,omitempty"`
	} `json:"predictions,omitempty"`
}
