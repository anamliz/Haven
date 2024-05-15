
package pollDataMysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/anamliz/Haven/internal/domains/pollData"
	_ "github.com/go-sql-driver/mysql" // MySQL driver import
)

var _ pollData.PollDataRepository = (*MysqlRepository)(nil)

type MysqlRepository struct {
	db *sql.DB
}

// Create a new mysql repository
func New(connectionString string) (*MysqlRepository, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	// Maximum Open Connections
	db.SetMaxOpenConns(10)
	// Idle Connection Timeout
	db.SetConnMaxIdleTime(5 * time.Second)
	// Connection Lifetime
	db.SetConnMaxLifetime(15 * time.Second)

	return &MysqlRepository{
		db: db,
	}, nil
}

// Save : saves live score record into db
func (mr *MysqlRepository) Save(ctx context.Context, y pollData.Accommodation) (int, error) {
	var d int
	rs, err := mr.db.Exec("INSERT accommodation SET id=?,name=?,description=?,price=?,imageurl=?, \n"+
		"comments =?,created=now(),update_at=now() ON DUPLICATE KEY UPDATE update_at=now()",
		y.ID, y.Name, y.Description, y.Price, y.Imageurl, y.Comments )

	if err != nil {
		return d, fmt.Errorf("Unable to save Accommodation : %v", err)
	}

	lastInsertedID, err := rs.LastInsertId()
	if err != nil {
		return d, fmt.Errorf("Unable to retrieve last id [primary key] : %v", err)
	}

	return int(lastInsertedID), nil
}

// Accommodation : query data used when querying for accommodation.
func (mr *MysqlRepository) Get(ctx context.Context) ([]pollData.Accommodation, error) {

	var gc []pollData.Accommodation

	statement := fmt.Sprintf("SELECT id,name,description,price, imageurl, \n" +
		"comments from accommodation ")
	raws, err := mr.db.Query(statement)
	if err != nil {
		return gc, err
	}
	for raws.Next() {
		var g pollData.Accommodation
		err := raws.Scan(&g.ID, &g.Name, &g.Description, &g.Price, &g.Imageurl,
			&g.Comments)
		if err != nil {
			return gc, err
		}
		gc = append(gc, g)
	}
	if raws.Err(); err != nil {
		return gc, err
	}
	raws.Close()

	return gc, nil
}
