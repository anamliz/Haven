package pollDataMysql

import (
	"context"
	"database/sql"
	"fmt"
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
		err := raws.Scan(&g.ID, &g.Name, &g.Description, &g.Price, &g.ImageURL, &g.Comments)
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



// FetchByID fetches data from the database by ID.
func (m *MysqlRepository) FetchByID(ctx context.Context, id int) (*pollDataTypes.Accommodation, error) {
	// Implementation for fetching data by ID from MySQL
	var data pollDataTypes.Accommodation
	query := `SELECT id, name, description, price, imageurl, comments, created_at, updated_at FROM accommodations WHERE id = ?`
	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&data.ID, &data.Name, &data.Description, &data.Price, &data.ImageURL, &data.Comments, &data.CreatedAt, &data.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &data, nil
}


// Update updates data in the database.
func (m *MysqlRepository) Update(ctx context.Context, id int, newData pollDataTypes.Accommodation) error {
	query := `UPDATE accommodations SET name=?, description=?, price=?, imageurl=?, comments=?, updated_at=? WHERE id=?`
	_, err := m.db.ExecContext(ctx, query, newData.Name, newData.Description, newData.Price, newData.ImageURL, newData.Comments, newData.UpdatedAt, id)
	return err
}



func (mr *MysqlRepository) Delete(ctx context.Context, id int) error {
    _, err := mr.db.Exec("DELETE FROM accommodation WHERE id = ?", id)
    if err != nil {
        return fmt.Errorf("unable to delete Accommodation: %v", err)
    }
    return nil
}
