/*
Package deps provides interfaces for all dependencies.
*/
package deps

import "encoding/json"

// Logger is implemented e. g. by the Logger of the standard library package
// log, but you can use a custom log if you like.
type Logger interface {
	Printf(format string, v ...any)
	Fatalf(format string, v ...any)
}

// Database provides methods to save to and retrieve objects from persistent
// datastore.
type Database interface {
	InsertCase(fields json.RawMessage) (int, error)
	UpdateCase(id int, fields json.RawMessage) error
	RetrieveCase(id int) (json.RawMessage, error)
	RetrieveCaseAll() (map[int]json.RawMessage, error)
}

// Environment provides the environment for this module.
type Environment interface {
	Host() string
	Port() string
	DBFilename() string
}
