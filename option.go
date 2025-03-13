package log

import "strings"

// Option is a container for optional properties that can be used for initializing the logging system.
type Option struct {
	level string
}

// WithLevel sets the logging level Option. The default logging level is info.
func WithLevel(level string) func(*Option) {
	return func(o *Option) {
		o.level = strings.TrimSpace(level)
	}
}
