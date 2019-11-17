package poster

// Poster posts a the given content somewhere..
type Poster interface {
	Post(content []byte) error
}
