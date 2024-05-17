package pollData

import (
	"fmt"
	"time"
)



// NewPollData creates a new AccommodationItem instance
func NewPollData(ID, Name, Description, Price, ImageURL, Comments string) (*Accommodation, error) {
	if ID == "" {
		return &Accommodation{}, fmt.Errorf("ID not set")
	}
	if Name == "" {
		return &Accommodation{}, fmt.Errorf("Name not set")
	}
	if Description == "" {
		return &Accommodation{}, fmt.Errorf("Description not set")
	}
	if Price == "" {
		return &Accommodation{}, fmt.Errorf("Price not set")
	}
	if ImageURL == "" {
		return &Accommodation{}, fmt.Errorf("ImageURL not set")
	}
	if Comments == "" {
		return &Accommodation{}, fmt.Errorf("Comments not set")
	}

	created := time.Now().Format("2006-01-02 15:04:05")
	updated := time.Now().Format("2006-01-02 15:04:05")

	// Final Object
	  return &Accommodation{
		Name:        Name,
		Description: Description,
		Price:       Price,
		ImageURL:    ImageURL,
		Comments:    Comments,
		CreatedAt:   created,
		UpdatedAt:   updated,
	}, nil
}
