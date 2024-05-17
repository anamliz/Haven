
package pollData

import "context"

type AccommodationFetcher interface {
    GetData(ctx context.Context) ([]Accommodation, error)
}


