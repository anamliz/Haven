package pollData

type Accommodation struct {
	Status string `json:"status"`
	Data   []AccommodationItem `json:"data"`
}


type AccommodationItem struct {

	Status string `json:"status"`
	Data   []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       string `json:"price"`
		Imageurl    string `json:"imageurl"`
		Comments    string `json:"comments"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	} `json:"data"`
}



