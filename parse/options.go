package parse

type config struct {
	DisallowUnknownFields bool
}

// Option modifies the parsing logic.
type Option func(cfg *config)

// WithDisallowUnknownFields is an option that will return an error if unknown field is supplied
func WithDisallowUnknownFields() Option {
	return func(cfg *config) {
		cfg.DisallowUnknownFields = true
	}
}
