
package pollData

import (
    "fmt"
    "time"
	"github.com/anamliz/Haven/internal/domains/pollDataTypes"
    
)

// NewPollData creates a new AccommodationItem instance
func NewPollData( Name, Description, Price, ImageURL, Comments string) (*pollDataTypes.Accommodation, error) {
    
    if Name == "" {
        return &pollDataTypes.Accommodation{}, fmt.Errorf("name not set")
    }
    if Description == "" {
        return &pollDataTypes.Accommodation{}, fmt.Errorf("description not set")
    }
    if Price == "" {
        return &pollDataTypes.Accommodation{}, fmt.Errorf("price not set")
    }
    if ImageURL == "" {
        return &pollDataTypes.Accommodation{}, fmt.Errorf("imageURL not set")
    }
    if Comments == "" {
        return &pollDataTypes.Accommodation{}, fmt.Errorf("comments not set")
    }

    created := time.Now().Format("2006-01-02 15:04:05")
    updated := time.Now().Format("2006-01-02 15:04:05")

    // Final Object
    return &pollDataTypes.Accommodation{
        Name:        Name,
        Description: Description,
        Price:       Price,
        ImageURL:    ImageURL,
        Comments:    Comments,
        CreatedAt:   created,
        UpdatedAt:   updated,
    }, nil
}

//


// UpdatePollData updates an existing Accommodation instance
func UpdatePollData(accommodation *pollDataTypes.Accommodation, Name, Description, Price, ImageURL, Comments string) error {
    if accommodation == nil {
        return fmt.Errorf("accommodation not found")
    }

    if Name != "" {
        accommodation.Name = Name
    }
    if Description != "" {
        accommodation.Description = Description
    }
    if Price != "" {
        accommodation.Price = Price
    }
    if ImageURL != "" {
        accommodation.ImageURL = ImageURL
    }
    if Comments != "" {
        accommodation.Comments = Comments
    }
    accommodation.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

    return nil
}

// DeletePollData deletes an Accommodation instance from a map of accommodations
func DeletePollData(accommodations map[string]*pollDataTypes.Accommodation, id string) error {
    if _, exists := accommodations[id]; !exists {
        return fmt.Errorf("accommodation with ID %s not found", id)
    }

    delete(accommodations, id)
    return nil
}

//

