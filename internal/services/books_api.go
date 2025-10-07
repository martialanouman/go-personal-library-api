package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type BookAPIClient struct {
	client  http.Client
	baseURL string
	token   string
	logger  *log.Logger
}

type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Rating struct {
	Average float64 `json:"average"`
}

type Identifiers struct {
	Isbn10 string `json:"isbn_10"`
	Isbn13 string `json:"isbn_13"`
}

type APIBook struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Image       string      `json:"image"`
	Authors     []Author    `json:"authors"`
	Rating      Rating      `json:"rating"`
	Identifiers Identifiers `json:"identifiers"`
	Description string      `json:"description"`
}

func NewBookAPIService(logger *log.Logger) (*BookAPIClient, error) {
	client := http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  false,
			DisableKeepAlives:   false,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}

			return nil
		},
	}

	token := os.Getenv("BIG_BOOK_API_TOKEN")
	if token == "" {
		return nil, errors.New("BIG_BOOK_API_TOKEN environment variable is not set")
	}
	baseURL := os.Getenv("BIG_BOOK_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.bigbookapi.com"
	}

	return &BookAPIClient{
		client:  client,
		baseURL: baseURL,
		token:   token,
		logger:  logger,
	}, nil
}

func (b *BookAPIClient) GetBookByBigBookId(ctx context.Context, bbId string) (*APIBook, error) {
	if bbId == "" {
		return nil, errors.New("bbId cannot be empty")
	}

	u, err := url.JoinPath(b.baseURL, bbId)
	if err != nil {
		return nil, fmt.Errorf("failed to join URL path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", b.token)
	req.Header.Set("Accept", "application/json")

	r, err := b.client.Do(req)
	if err != nil {
		b.logger.Printf("request failed with error: %v", err)
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		b.logger.Printf("request failed with status code: %d", r.StatusCode)
		switch r.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("book with ID %s not found", bbId)
		case http.StatusUnauthorized:
			return nil, errors.New("unauthorized: invalid API token")
		case http.StatusTooManyRequests:
			return nil, errors.New("rate limit exceeded: too many requests")
		default:
			return nil, fmt.Errorf("request failed with status code: %d", r.StatusCode)
		}
	}

	var book APIBook

	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		b.logger.Printf("failed to decode response body: %v", err)
		return nil, err
	}

	return &book, nil
}
