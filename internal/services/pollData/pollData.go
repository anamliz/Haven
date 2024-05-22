package pollData

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anamliz/Haven/internal/domains/client/polldata"
	"github.com/anamliz/Haven/internal/domains/pollDataTypes"
	"github.com/anamliz/Haven/internal/domains/pollData"
	pollDataMysql "github.com/anamliz/Haven/internal/domains/pollData/pollDataMysql"
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

// Post is used to create new poll data and save it to the database.
func (s *PollDataService) Post(ctx context.Context, newData pollDataTypes.Accommodation) (int, error) {
    log.Printf("Starting creation of new accommodation data")

    // Save the new data into the database
    lastInsertedID, err := s.pollMysql.Save(ctx, newData)
    if err != nil {
        log.Printf("Error saving new data to MySQL: %v", err)
        return 0, err
    }

    log.Printf("New data created successfully with ID: %d", lastInsertedID)
    return lastInsertedID, nil
}


func (s *PollDataService) Update(ctx context.Context, id int, newData pollDataTypes.Accommodation) error {
    log.Printf("Starting update for accommodation with ID: %d", id)

    existingData, err := s.pollMysql.FetchByID(ctx, id)
    if err != nil {
        log.Printf("Failed to fetch existing data for ID: %d, error: %v", id, err)
        return err
    }

    existingData.Name = newData.Name
    existingData.Description = newData.Description
    existingData.Price = newData.Price
    existingData.ImageURL = newData.ImageURL
    existingData.Comments = newData.Comments

    err = s.pollMysql.Update(ctx, id, *existingData)
    if err != nil {
        log.Printf("Update failed for accommodation with ID: %d, error: %v", id, err)
        return err
    }

    log.Printf("Update successful for accommodation with ID: %d", id)
    return nil
}

func (s *PollDataService) Delete(ctx context.Context, id int) error {
    log.Printf("Starting delete for accommodation with ID: %d", id)

    err := s.pollMysql.Delete(ctx, id)
    if err != nil {
        log.Printf("Delete failed for accommodation with ID: %d, error: %v", id, err)
        return err
    }

    log.Printf("Delete successful for accommodation with ID: %d", id)
    return nil
}
