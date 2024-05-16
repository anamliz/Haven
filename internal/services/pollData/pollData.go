package pollData

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anamliz/Haven/internal/domains/client/polldata"
	"github.com/anamliz/Haven/internal/domains/pollData"
	pollDataMysql "github.com/anamliz/Haven/internal/domains/pollData/pollDataMysql"
)

// PollDataServicesConfiguration is an alias for a function that will take in a pointer to a PollDataService and modify it
type PollDataServicesConfiguration func(os *PollDataService) error

// PollDataService is an implementation of the PollDataService
type PollDataService struct {
	pollMysql pollData.PollDataRepository
}

// Data represents data retrieved from an external API
type Data struct {
	RawData string
}

// NewPollDataService instantiates every connection needed to run the PollData service
func NewPollDataService(cfgs ...PollDataServicesConfiguration) (*PollDataService, error) {
	os := &PollDataService{}
	for _, cfg := range cfgs {
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}
	return os, nil
}

// WithMysqlPollDataRepository instantiates MySQL to connect to the matches interface
func WithMysqlPollDataRepository(connectionString string) PollDataServicesConfiguration {
	return func(os *PollDataService) error {
		d, err := pollDataMysql.New(connectionString)
		if err != nil {
			return err
		}
		os.pollMysql = d
		return nil
	}
}

// PollData processes live scores
func (s *PollDataService) PollData(ctx context.Context, pollDataEndPoint string, timeouts time.Duration, client *http.Client) error {
	d, err := polldata.New(pollDataEndPoint, timeouts, client)
	if err != nil {
		return err
	}

	data, err := d.GetData(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}

	for _, d := range data.Data {
		log.Printf("*** ID: %s | Name: %s ", d.ID, d.Name)

		newData := &pollData.AccommodationItem{
			ID:          d.ID,
			Name:        d.Name,
			Description: d.Description,
			Price:       d.Price,
			ImageURL:    d.ImageURL, // Assuming ImageURL is the correct field
			Comments:    d.Comments,
		}

		_, err = s.pollMysql.Save(ctx, *newData)
		if err != nil {
			log.Printf("Error saving data: %v", err)
		}
	}
	return nil
}
