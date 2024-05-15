package pollData

import "context"

type DataFetcher interface {
	GetData(ctx context.Context) ([]Accommodation, error)
}

