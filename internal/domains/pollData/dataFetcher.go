/*package pollData

import "context"

type AccommodationFetcher interface {
	GetData(ctx context.Context) ([]Accommodation, error)
}*/


package pollData

import (
    "context"
    "github.com/anamliz/Haven/internal/domains/pollDataTypes"
)

type AccommodationFetcher interface {
    GetData(ctx context.Context) ([]pollDataTypes.Accommodation, error)
}