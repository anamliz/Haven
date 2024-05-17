

package pollData

import "context"

type PollDataRepository interface {
	Save(context.Context, Accommodation) (int, error)
	Get(ctx context.Context) ([]Accommodation , error)
}

