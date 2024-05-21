
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
