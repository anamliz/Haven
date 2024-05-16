package pollData

import (
	"fmt"
	"time"
)



// NewPollData creates a new AccommodationItem instance
func NewPollData(ID, Name, Description, Price, ImageURL, Comments string) (*AccommodationItem, error) {
	if ID == "" {
		return nil, fmt.Errorf("ID not set")
	}
	if Name == "" {
		return nil, fmt.Errorf("Name not set")
	}
	if Description == "" {
		return nil, fmt.Errorf("Description not set")
	}
	if Price == "" {
		return nil, fmt.Errorf("Price not set")
	}
	if ImageURL == "" {
		return nil, fmt.Errorf("ImageURL not set")
	}
	if Comments == "" {
		return nil, fmt.Errorf("Comments not set")
	}

	created := time.Now().Format("2006-01-02 15:04:05")
	updated := time.Now().Format("2006-01-02 15:04:05")

	// Create and return the AccommodationItem instance
	item := &AccommodationItem{
		ID:          ID,
		Name:        Name,
		Description: Description,
		Price:       Price,
		ImageURL:    ImageURL,
		Comments:    Comments,
		CreatedAt:   created,
		UpdatedAt:   updated,
	}

	return item, nil
}
