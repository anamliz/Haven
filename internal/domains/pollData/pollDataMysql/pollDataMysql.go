package pollDataMysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/anamliz/Haven/internal/domains/pollData"
	"github.com/anamliz/Haven/internal/domains/pollDataTypes"
)

var _ pollData.PollDataRepository = (*MysqlRepository)(nil)

type MysqlRepository struct {
	db *sql.DB
}

// New creates a new MySQL repository
func New(connectionString string) (*MysqlRepository, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(5 * time.Second)
	db.SetConnMaxLifetime(15 * time.Second)

	return &MysqlRepository{
		db: db,
	}, nil
}

func (mr *MysqlRepository) Save(ctx context.Context, y pollDataTypes.Accommodation) (int, error) {
	var d int
	rs, err := mr.db.Exec("INSERT INTO accommodation (name, description, price, imageurl, comments, created_at, updated_at) "+
		"VALUES (?, ?, ?, ?, ?, NOW(), NOW()) ON DUPLICATE KEY UPDATE updated_at=NOW()",
		y.Name, y.Description, y.Price, y.ImageURL, y.Comments)
	if err != nil {
		return d, fmt.Errorf("unable to save Accommodation: %v", err)
	}

	lastInsertedID, err := rs.LastInsertId()
	if err != nil {
		return d, fmt.Errorf("unable to retrieve last id [primary key]: %v", err)
	}

	return int(lastInsertedID), nil
}

func (mr *MysqlRepository) Get(ctx context.Context) ([]pollDataTypes.Accommodation, error) {
	var gc []pollDataTypes.Accommodation

	statement := "SELECT id, name, description, price, imageurl, comments FROM accommodation"
	raws, err := mr.db.Query(statement)
	if err != nil {
		return gc, err
	}
	for raws.Next() {
		var g pollDataTypes.Accommodation
		err := raws.Scan(&g.ID, &g.Name, &g.Description, &g.Price, &g.ImageURL, &g.Comments, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return gc, err
		}
		gc = append(gc, g)
	}
	if err := raws.Err(); err != nil {
		return gc, err
	}
	raws.Close()

	return gc, nil
}

// FetchByID retrieves poll data from the database based on the provided ID.
func (mr *MysqlRepository) FetchByID(ctx context.Context, id int) (*pollDataTypes.Accommodation, error) {
	// Database query to fetch data based on ID
	
	row := mr.db.QueryRowContext(ctx, "SELECT * FROM accommodation WHERE id = ?", id)

	// Scan the query result into a struct
	var data pollDataTypes.Accommodation
	err := row.Scan(&data.ID, &data.Name, &data.Description, &data.Price, &data.ImageURL, &data.Comments, &data.CreatedAt, &data.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No data found for ID: %d", id)
			return nil, nil // Return nil if no data found, without error
		}
		log.Printf("Error fetching data for ID: %d, error: %v", id, err)
		return nil, err
	}

	return &data, nil
}

// Update updates the poll data in the database based on the provided ID.
func (mr *MysqlRepository) Update(ctx context.Context, id int, newData pollDataTypes.Accommodation) error {
	// Database query to update data based on ID
	// Example query: UPDATE accommodation SET name = ?, description = ?, price = ?, imageurl = ?, comments = ? WHERE id = ?
	_, err := mr.db.ExecContext(ctx, "UPDATE accommodation SET name = ?, description = ?, price = ?, imageurl = ?, comments = ? WHERE id = ?",
		newData.Name, newData.Description, newData.Price, newData.ImageURL, newData.Comments, id)
	if err != nil {
		log.Printf("Error updating data for ID: %d, error: %v", id, err)
		return err
	}

	log.Printf("Data updated successfully for ID: %d", id)
	return nil
}

func (mr *MysqlRepository) Delete(ctx context.Context, id int) error {
	// Database query to delete data based on ID
	// Example query: DELETE FROM accommodation WHERE id = ?
	_, err := mr.db.ExecContext(ctx, "DELETE FROM accommodation WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting data for ID: %d, error: %v", id, err)
		return err
	}

	log.Printf("Data deleted successfully for ID: %d", id)
	return nil
}
