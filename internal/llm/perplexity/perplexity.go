package perplexity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/transport"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client

	token string
}

func NewClient(token string) (*Client, error) {
	baseURL := "https://api.perplexity.ai"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("perplexity: failed to parse base URL: %w", err)
	}

	httpCli := &http.Client{
		Timeout: 60 * time.Second,
		Transport: transport.Chain(
			http.DefaultTransport,
			transport.SetHeader("Authorization", "Bearer "+token),
			transport.SetHeader("Content-Type", "application/json"),
		),
	}

	return &Client{
		baseURL:    u,
		httpClient: httpCli,
		token:      token,
	}, nil
}

type Request interface {
	Path() string
	Method() string
}

func do[T interface{ Request }, U any](httpClient *http.Client, baseURL *url.URL, req T) (*U, error) {
	u := *baseURL
	u.Path = req.Path()
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("perplexity: failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(req.Method(), u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("perplexity: failed to create request: %w", err)
	}

	var data U
	res, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("perplexity: failed to do request: %w", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("perplexity: failed to decode response: %w", err)
	}

	return &data, nil
}

type ChatCompletionMessage struct {
	Content string `json:"content"`
	Role    Role   `json:"role"`
}

type ChatCompletionRequest struct {
	Model                  string                  `json:"model"`
	Messages               []ChatCompletionMessage `json:"messages"`
	MaxTokens              int                     `json:"max_tokens,omitempty"`
	Temperature            float64                 `json:"temperature,omitempty"`
	TopP                   float64                 `json:"top_p,omitempty"`
	ReturnCitations        bool                    `json:"return_citations,omitempty"`
	SearchDomainFilter     []string                `json:"search_domain_filter,omitempty"`
	ReturnImages           bool                    `json:"return_images,omitempty"`
	ReturnRelatedQuestions bool                    `json:"return_related_questions,omitempty"`
	SearchRecencyFilter    string                  `json:"search_recency_filter,omitempty"`
	TopK                   int                     `json:"top_k,omitempty"`
	Stream                 bool                    `json:"stream,omitempty"`
	PresencePenalty        float64                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty       float64                 `json:"frequency_penalty,omitempty"`
}

func NewChatCompletionRequest() ChatCompletionRequest {
	return ChatCompletionRequest{
		Model: "llama-3.1-sonar-small-128k-online",
		// Messages:
		// MaxTokens:
		Temperature:     0.2,
		TopP:            0.9,
		ReturnCitations: false,
		// SearchDomainFilter:
		ReturnImages:           false,
		ReturnRelatedQuestions: false,
		// SearchRecencyFilter:
		TopK:             0,
		Stream:           false, // TODO set to true
		PresencePenalty:  0,
		FrequencyPenalty: 1,
	}
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index        int                   `json:"index"`
		FinishReason string                `json:"finish_reason"`
		Message      ChatCompletionMessage `json:"message"`
		// delta
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (c ChatCompletionRequest) Path() string {
	return "/chat/completions"
}

func (c ChatCompletionRequest) Method() string {
	return http.MethodPost
}

func (c *Client) ChatCompletions(r ChatCompletionRequest) (*ChatCompletionResponse, error) {
	res, err := do[ChatCompletionRequest, ChatCompletionResponse](c.httpClient, c.baseURL, r)
	if err != nil {
		return nil, fmt.Errorf("perplexity: failed to do chat completions: %w", err)
	}

	return res, nil
}
