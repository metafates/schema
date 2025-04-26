package parse

func defaultConfig() config {
	return config{
		DisallowUnknownFields: false,
		RenameFunc:            func(s string) string { return s },
	}
}

type config struct {
	DisallowUnknownFields bool
	RenameFunc            RenameFunc
}

type RenameFunc func(string) string

// Option modifies the parsing logic.
type Option func(cfg *config)

// WithDisallowUnknownFields is an option that will return an error if unknown field is supplied
func WithDisallowUnknownFields() Option {
	return func(cfg *config) {
		cfg.DisallowUnknownFields = true
	}
}

// WithRenameFunc is an option that will rename src fields/keys during parsing before matching with dst fields
func WithRenameFunc(f RenameFunc) Option {
	return func(cfg *config) {
		cfg.RenameFunc = f
	}
}
