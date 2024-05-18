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
	"github.com/anamliz/Haven/internal/domains/pollDataTypes"
)

// Compile time interface assertion.
var _ pollData.AccommodationFetcher = (*PollDataClient)(nil)

// PollDataClient implements the DataFetcher interface.
type PollDataClient struct {
	pollDataEndPoint string
	timeouts         time.Duration
	client           *http.Client
}

// New initializes a new instance of PollDataClient.
func New(pollDataEndPoint string, timeouts time.Duration, client *http.Client) (*PollDataClient, error) {
	pollDataURL, err := url.Parse(pollDataEndPoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse poll data endpoint: %w", err)
	}

	if timeouts <= 0 {
		return nil, fmt.Errorf("timeout not set")
	}

	c := &PollDataClient{
		pollDataEndPoint: pollDataURL.String(),
		timeouts:         timeouts,
		client:           client,
	}
	if c.client == nil {
		c.client = defaultHTTPClient
	}

	return c, nil
}

// GetData fetches data using PollDataClient.
func (s *PollDataClient) GetData(ctx context.Context) ([]pollDataTypes.Accommodation, error) {
	var pollDataURL = s.pollDataEndPoint
	log.Printf("Calling... : %s", pollDataURL)

	response, err := s.client.Get(pollDataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call poll data API: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if response.StatusCode == http.StatusOK {
		return parseApiData(responseBody)
	}

	return nil, fmt.Errorf("failed to get data: status: %d", response.StatusCode)
}

// parseApiData parses the JSON response into Accommodation.
func parseApiData(content []byte) ([]pollDataTypes.Accommodation, error) {
	var g []pollDataTypes.Accommodation

	var s pollDataTypes.Response

	err := json.Unmarshal(content, &s)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	log.Printf("status: %s", s.Status)

	for _, i := range s.Data {
		log.Printf("ID: %s | Name: %s", i.ID, i.Name)

		d := pollDataTypes.Accommodation{

			Name:        i.Name,
			Description: i.Description,
			Price:       i.Price,
			ImageURL:    i.ImageURL,
			Comments:    i.Comments,
		}
		g = append(g, d)
	}

	if len(g) == 0 {
		return g, fmt.Errorf("accommodations list is empty")
	}

	return g, nil
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
