package pollData

import (
	"fmt"
	"time"
)

func NewPollData(ID, Name , Description, Price,Imageurl,Comments string) (*Accommodation, error) {

	if ID == "" {
		return &Accommodation{}, fmt.Errorf("ID not set")
	}

	if Name == "" {
		return &Accommodation{}, fmt.Errorf("Name not set")
	}

	if Description == "" {
		return &Accommodation{}, fmt.Errorf("Description not set")
	}

    if Price== "" {
		return &Accommodation{}, fmt.Errorf("Price not set")
	}

    if Imageurl == "" {
		return &Accommodation{}, fmt.Errorf("Imageurl not set")
	}

    if Comments == "" {
		return &Accommodation{}, fmt.Errorf("Comments not set")
	}


	created_at := time.Now().Format("2006-01-02 15:04:05")

	// Final Object
	return &Accommodation{
		ID :  id,
		Name: name,
		Description: description,
		Price : price,
		Imageurl: imageurl,
        Comments: comments,
		Created:  created_at,
        UpdatedAt: updated_at,
	}, nil

}
