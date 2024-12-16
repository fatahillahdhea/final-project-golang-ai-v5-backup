package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"a21hc3NpZ25tZW50/model"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	if len(table) == 0 {
		return "", errors.New("table cannot be empty")
	}

	payload := map[string]interface{}{
		"table": table,
		"query": query,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Cells []string `json:"cells"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Cells) == 0 {
		return "", errors.New("empty result received from API")
	}

	return result.Cells[0], nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	type Messages struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type RequestBody struct {
		Messages []Messages `json:"messages"`
	}

	payload := RequestBody{
		Messages: []Messages{
			{
				Role:    "user",
				Content: query,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Validasi Status Code
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return model.ChatResponse{}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Struct untuk menangani response
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type Logprobs struct {
		// Jika struktur logprobs lebih kompleks, tambahkan di sini
	}

	type Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
	type Choice struct {
		Index        int       `json:"index"`
		Message      Message   `json:"message"`
		Logprobs     *Logprobs `json:"logprobs"` // Nullable, jadi pointer
		FinishReason string    `json:"finish_reason"`
	}

	type ChatCompletionResponse struct {
		Object            string   `json:"object"`
		ID                string   `json:"id"`
		Created           int64    `json:"created"`
		Model             string   `json:"model"`
		SystemFingerprint string   `json:"system_fingerprint"`
		Choices           []Choice `json:"choices"`
		Usage             Usage    `json:"usage"`
	}

	var result ChatCompletionResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Validasi apakah array tidak kosong
	// if len(result) == 0 || result[0].GeneratedText == "" {
	// 	return model.ChatResponse{}, errors.New("empty response received from API")
	// }

	// Mengambil teks hasil pertama dari array
	return model.ChatResponse{GeneratedText: result.Choices[0].Message.Content}, nil
}
