
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

	"github.com/anamliz/Haven/internal/domains/pollData"
)

// Compile time interface assertion.
var _ pollData.DataFetcher = (*PollDataClient)(nil)

type PollDataClient struct {
	accommodationEndPoint string
	timeouts              time.Duration
	client                *http.Client
}

// New initializes a new instance of Accommodation Client.
func New(accommodationEndPoint string, timeouts time.Duration, client *http.Client) (*PollDataClient, error) {

	accommodationURL, err := url.Parse(accommodationEndPoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse accommodation endpoint: %w", err)
	}

	if timeouts <= 0 {
		return nil, fmt.Errorf("Timeout not set")
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

func (s *PollDataClient) GetData(ctx context.Context) ([]pollData.AccommodationItem, error) {

	var accommodationURL = fmt.Sprintf("%s", s.accommodationEndPoint)
	log.Printf("Calling... : %s", accommodationURL)

	response, err := s.client.Get(accommodationURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to call accommodation API: %v", err)
	}

	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if response.StatusCode == http.StatusOK {
		return parseAccommodationData(responseBody)
	}

	return nil, fmt.Errorf("failed to get data : status: %d, error body: %s", response.StatusCode, responseBody)

}

func parseAccommodationData(content []byte) ([]pollData.AccommodationItem, error) {

	var accommodations []pollData.AccommodationItem

	var response pollData.Accommodation
	err := json.Unmarshal(content, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON : %v", err)
	}

	log.Printf("Status: %s", response.Status)

	for _, item := range response.Data {

		log.Printf("Accommodation ID: %s | Name: %s", item.ID, item.Name)

		accommodation := pollData.AccommodationItem{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			ImageURL:    item.ImageURL,
			Comments:    item.Comments,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}

		accommodations = append(accommodations, accommodation)

	}

	if len(accommodations) == 0 {
		return accommodations, fmt.Errorf("accommodations list is empty")
	}

	return accommodations, nil
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
