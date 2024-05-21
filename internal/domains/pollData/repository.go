
package pollData

import (
    "context"
	"github.com/anamliz/Haven/internal/domains/pollDataTypes"
    
)

type PollDataRepository interface {
    Save(context.Context, pollDataTypes.Accommodation) (int, error)
    Get(ctx context.Context) ([]pollDataTypes.Accommodation, error)
    FetchByID(ctx context.Context, id int) (*pollDataTypes.Accommodation, error)
    Delete(ctx context.Context, id int) error
    Update(ctx context.Context, id int, newData pollDataTypes.Accommodation) error
}


