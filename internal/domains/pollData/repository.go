
package pollData

import (
    "context"
	"github.com/anamliz/Haven/internal/domains/pollDataTypes"
    
)

type PollDataRepository interface {
    Save(context.Context, pollDataTypes.Accommodation) (int, error)
    Get(ctx context.Context) ([]pollDataTypes.Accommodation, error)
}