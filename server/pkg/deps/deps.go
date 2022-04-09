/*
Package deps provides interfaces for all dependencies.
*/
package deps

// Logger is implemented e. g. by the Logger of the standard library package
// log, but you can use a custom log if you like.
type Logger interface {
	Printf(format string, v ...any)
	Fatalf(format string, v ...any)
}

// GetEnvFunc is a function that retrieves environment variables. Use os.Getenv
// in production.
type GetEnvFunc func(key string) string
