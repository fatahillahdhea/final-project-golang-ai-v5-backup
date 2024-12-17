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

	// Mengambil jawaban yang relevan, bukan hanya yang terakhir
	return result.Cells[0], nil // Atau logika lain untuk memilih jawaban yang tepat
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	// Membuat request body yang lebih sederhana
	requestBody := map[string]interface{}{
		"inputs": query,
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	modelUrl := "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct"
	req, err := http.NewRequest("POST", modelUrl, bytes.NewReader(reqBody))
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

	// Baca body respons
	var response []model.ChatResponse // array, karena API mengembalikan array
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response) == 0 {
		return model.ChatResponse{}, errors.New("empty response received from API")
	}

	return response[0], nil
}