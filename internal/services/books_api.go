package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type BookAPIClient struct {
	client  http.Client
	baseURL string
	token   string
	logger  *log.Logger
}

type APIBook struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Image   string `json:"image"`
	Authors []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"authors"`
	Rating struct {
		Average float64 `json:"average"`
	} `json:"rating"`
	Identifiers struct {
		Isbn10 string `json:"isbn_10"`
		Isbn13 string `json:"isbn_13"`
	} `json:"identifiers"`
	Description string `json:"description"`
}

func NewBookAPIService(logger *log.Logger) *BookAPIClient {
	client := http.Client{
		Timeout: 30 * time.Second,
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

	token := os.Getenv("BOOK_API_TOKEN")
	if token == "" {
		panic("BOOK_API_TOKEN is not set")
	}

	return &BookAPIClient{
		client:  client,
		baseURL: "https://api.bigbookapi.com",
		token:   token,
		logger:  logger,
	}
}

func (b *BookAPIClient) GetBookByBigBookId(bbId string) (*APIBook, error) {
	url := b.baseURL + "/" + bbId

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
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
		resBody, err := io.ReadAll(r.Body)
		if err != nil {
			b.logger.Printf("failed to read response body: %v", err)
		}
		b.logger.Printf("request failed with status code: %d, response: %s", r.StatusCode, string(resBody))
		return nil, fmt.Errorf("request failed with status code: %d", r.StatusCode)
	}

	var book *APIBook

	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		b.logger.Printf("failed to decode response body: %v", err)
		return nil, err
	}

	if book == nil {
		return nil, nil
	}

	return book, nil
}
