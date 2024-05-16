
package pollData

// AccommodationItem represents a single accommodation item
type AccommodationItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	ImageURL    string `json:"imageurl"`
	Comments    string `json:"comments"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// AccommodationResponse represents the JSON response structure for accommodation data
type AccommodationResponse struct {
	Status string             `json:"status"`
	Data   []AccommodationItem `json:"data"`
}
