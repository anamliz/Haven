package pollDataTypes

type Accommodation struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	ImageURL    string `json:"imageurl"`
	Comments    string `json:"comments"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Response struct {
	Status string          `json:"status"`
	Data   []Accommodation `json:"data"`
}
