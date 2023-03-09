package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Api interface {
	Completions(request CompletionRequest) (CompletionResponse, error)
	ChatCompletions(request ChatCompletionRequest) (ChatCompletionResponse, error)
}

type api struct {
	client *http.Client
	url    string
	key    string
}

func (a api) Completions(request CompletionRequest) (CompletionResponse, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return CompletionResponse{}, err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/completions", a.url),
		bytes.NewBuffer(requestBytes),
	)
	if err != nil {
		return CompletionResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.key))

	resp, err := a.client.Do(req)
	if err != nil {
		return CompletionResponse{}, err
	}
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return CompletionResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return CompletionResponse{}, fmt.Errorf("不正なAPIリクエストです")
	}
	var response CompletionResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return CompletionResponse{}, err
	}
	return response, nil
}

func (a api) ChatCompletions(request ChatCompletionRequest) (ChatCompletionResponse, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/chat/completions", a.url),
		bytes.NewBuffer(requestBytes),
	)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.key))

	resp, err := a.client.Do(req)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return ChatCompletionResponse{}, fmt.Errorf("不正なAPIリクエストです")
	}
	var response ChatCompletionResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	return response, nil
}

func NewApi(client *http.Client, url, key string) Api {
	return &api{client, url, key}
}

type CompletionRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int64  `json:"max_tokens"`
}

type CompletionResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type ChatCompletionRequest struct {
	Model     string                  `json:"model"`
	Messages  []ChatCompletionMessage `json:"messages"`
	MaxTokens int64                   `json:"max_tokens"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
