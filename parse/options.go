package parse

type config struct {
	DisallowUnknownFields bool
}

type Option func(cfg *config)

func WithDisallowUnknownFields() Option {
	return func(cfg *config) {
		cfg.DisallowUnknownFields = true
	}
}
