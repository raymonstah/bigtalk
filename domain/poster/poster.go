package poster

import "context"

// Poster posts a the given content somewhere..
type Poster interface {
	Post(ctx context.Context, content []byte) error
}
