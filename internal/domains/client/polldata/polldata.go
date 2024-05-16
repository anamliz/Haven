package polldata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

// DataFetcher defines the interface for fetching data.
type DataFetcher interface {
	GetData(ctx context.Context) ([]AccommodationItem, error)
}

// Compile time interface assertion.
var _ DataFetcher = (*PollDataClient)(nil)

// PollDataClient implements the DataFetcher interface.
type PollDataClient struct {
	accommodationEndPoint string
	timeouts              time.Duration
	client                *http.Client
}

// New initializes a new instance of PollDataClient.
func New(accommodationEndPoint string, timeouts time.Duration, client *http.Client) (*PollDataClient, error) {
	accommodationURL, err := url.Parse(accommodationEndPoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse accommodation endpoint: %w", err)
	}

	if timeouts <= 0 {
		return nil, fmt.Errorf("timeout not set")
	}

	c := &PollDataClient{
		accommodationEndPoint: accommodationURL.String(),
		timeouts:              timeouts,
		client:                client,
	}
	if c.client == nil {
		c.client = defaultHTTPClient
	}

	return c, nil
}

// GetData fetches data using PollDataClient.
func (s *PollDataClient) GetData(ctx context.Context) ([]AccommodationItem, error) {
	var accommodationURL = fmt.Sprintf("%s", s.accommodationEndPoint)
	log.Printf("Calling... : %s", accommodationURL)

	response, err := s.client.Get(accommodationURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call accommodation API: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if response.StatusCode == http.StatusOK {
		return parseAccommodationData(responseBody)
	}

	return nil, fmt.Errorf("failed to get data: status: %d", response.StatusCode)
}

// parseAccommodationData parses the JSON response into AccommodationItem.
func parseAccommodationData(content []byte) ([]AccommodationItem, error) {
	var accommodations []AccommodationItem
	var response struct {
		Status string
		Data   []AccommodationItem
	}

	err := json.Unmarshal(content, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	log.Printf("Status: %s", response.Status)

	for _, item := range response.Data {
		log.Printf("Accommodation ID: %s | Name: %s", item.ID, item.Name)
		accommodations = append(accommodations, item)
	}

	if len(accommodations) == 0 {
		return nil, fmt.Errorf("accommodations list is empty")
	}

	return accommodations, nil
}

// AccommodationItem represents an item of accommodation data.
type AccommodationItem struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url"`
	Comments    []string  `json:"comments"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var defaultHTTPClient = &http.Client{
	Timeout: time.Second * 15,
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Second * 15,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	},
}
