package pollData

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/anamliz/Haven/internal/domains/pollData/pollDataTypes"
)

// TestUpdateAccommodation tests the Update method of PollDataService.
func TestUpdateAccommodation(t *testing.T) {
	// Set up the MySQL connection
	connectionString := "your-mysql-connection-string"
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Create the service
	service, err := NewPollDataService(WithMysqlPollDataRepository(connectionString))
	if err != nil {
		log.Fatalf("Failed to create PollDataService: %v", err)
	}

	ctx := context.Background()
	id := 1 // Assuming this ID exists in your database
	newData := pollDataTypes.Accommodation{
		Name:        "Updated Name",
		Description: "Updated Description",
		Price:       "123.45",
		ImageURL:    "http://newimage.url",
		Comments:    "Updated comments",
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	err = service.Update(ctx, id, newData)
	if err != nil {
		log.Printf("Update failed: %v", err)
	} else {
		log.Printf("Update succeeded for ID %d", id)
	}
}

// TestDeleteAccommodation tests the DeleteAccommodation method of PollDataService.
func TestDeleteAccommodation(t *testing.T) {
	// Set up the MySQL connection
	connectionString := "your-mysql-connection-string"
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Create the service
	service, err := NewPollDataService(WithMysqlPollDataRepository(connectionString))
	if err != nil {
		log.Fatalf("Failed to create PollDataService: %v", err)
	}

	ctx := context.Background()
	id := 1 // Assuming this ID exists in your database

	err = service.DeleteAccommodation(ctx, id)
	if err != nil {
		log.Printf("Delete failed: %v", err)
	} else {
		log.Printf("Delete succeeded for ID %d", id)
	}
}
