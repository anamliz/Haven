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
	"github.com/anamliz/Haven/internal/domains/pollDataTypes"
)

// PollDataServicesConfiguration is an alias for a function that will take in a pointer to a PollDataService and modify it
type PollDataServicesConfiguration func(os *PollDataService) error

// PollDataService is an implementation of the PollDataService
type PollDataService struct {
	pollMysql pollData.PollDataRepository
}

type Data struct {
	RawData string
}

// NewPollDataService : instantiate every connection we need to run season service
func NewPollDataService(cfgs ...PollDataServicesConfiguration) (*PollDataService, error) {
	// Create the PollDataService
	os := &PollDataService{}
	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		// Pass the service into the configuration function
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}
	return os, nil
}

// WithMysqlPollDataRepository : instantiates mysql to connect to matches interface
func WithMysqlPollDataRepository(connectionString string) PollDataServicesConfiguration {
	return func(os *PollDataService) error {
		// Create PollData repository
		d, err := pollDataMysql.New(connectionString)
		if err != nil {
			return err
		}
		os.pollMysql = d
		return nil
	}
}

// PollData : processes accommodation
func (s *PollDataService) PollData(ctx context.Context, pollDataEndPoint string, timeouts time.Duration, client *http.Client) error {
	// Poll Data from external API
	d, err := polldata.New(pollDataEndPoint, timeouts, client)
	if err != nil {
		return err
	}

	data, err := d.GetData(ctx)
	if err != nil {
		return fmt.Errorf("unable to fetch data | %v", err)
	}

	for _, d := range data {
		log.Printf("*** Name: %s | Description: %s | Price : %s |ImageURL : %s", d.Name, d.Description, d.Price, d.ImageURL)

		// Save into database
		data, err := pollData.NewPollData(d.Name, d.Description, d.Price, d.ImageURL, d.Comments)
		if err != nil {
			log.Printf("Err : %s", err)
		} else {

			_, err = s.pollMysql.Save(ctx, *data)
			if err != nil {
				//log.Printf("Err : %s", err)
				log.Printf("Error saving data to MySQL: %v", err)
			} else {
				log.Printf("Data saved successfully to MySQL")
			}
		}

	}

	return nil

}

// Update updates an accommodation in the database.
func (s *PollDataService) Update(ctx context.Context, id int, newData pollDataTypes.Accommodation) error {
	// Fetch the existing data from the database based on the provided ID.
	existingData, err := s.pollMysql.FetchByID(ctx, id)
	if err != nil {
		return err
	}

	// Update the fields of the existing data with the new data.
	existingData.Name = newData.Name
	existingData.Description = newData.Description
	existingData.Price = newData.Price
	existingData.ImageURL = newData.ImageURL
	existingData.Comments = newData.Comments
	existingData.UpdatedAt = time.Now().Format(time.RFC3339)

	// Save the updated data back to the database.
	err = s.pollMysql.Update(ctx, id, *existingData)
	if err != nil {
		return err
	}

	// Log the successful update.
	log.Printf("Data with ID %d updated successfully: %+v", id, existingData)

	return nil
}

// DeleteAccommodation deletes accommodation by ID.
func (s *PollDataService) DeleteAccommodation(ctx context.Context, id int) error {
	err := s.pollMysql.Delete(ctx, id)
	if err != nil {
		log.Printf("Error deleting accommodation with ID %d: %v", id, err)
		return fmt.Errorf("failed to delete accommodation: %v", err)
	}

	log.Printf("Accommodation with ID %d deleted successfully", id)
	return nil
}
