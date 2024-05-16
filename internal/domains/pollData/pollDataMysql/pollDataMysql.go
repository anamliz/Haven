

package pollDataMysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/anamliz/Haven/internal/domains/pollData"
)

type MysqlRepository struct {
	db *sql.DB
}

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

func (mr *MysqlRepository) Save(ctx context.Context, y pollData.AccommodationItem) (int, error) {
	var d int
	rs, err := mr.db.Exec("INSERT accommodation SET id=?,name=?,description=?,price=?,imageurl=?, \n"+
		"comments =?,created=now(),update_at=now() ON DUPLICATE KEY UPDATE update_at=now()",
		y.ID, y.Name, y.Description, y.Price, y.ImageURL, y.Comments)

	if err != nil {
		return d, fmt.Errorf("Unable to save Accommodation : %v", err)
	}

	lastInsertedID, err := rs.LastInsertId()
	if err != nil {
		return d, fmt.Errorf("Unable to retrieve last id [primary key] : %v", err)
	}

	return int(lastInsertedID), nil
}

func (mr *MysqlRepository) Get(ctx context.Context) ([]pollData.AccommodationItem, error) {
	var gc []pollData.AccommodationItem

	statement := fmt.Sprintf("SELECT id,name,description,price, imageurl, \n" +
		"comments from accommodation ")
	raws, err := mr.db.Query(statement)
	if err != nil {
		return gc, err
	}
	for raws.Next() {
		var g pollData.AccommodationItem
		err := raws.Scan(&g.ID, &g.Name, &g.Description, &g.Price, &g.ImageURL,
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
