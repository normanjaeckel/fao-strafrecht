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
	Insert(name string, data json.RawMessage) (int, error)
	Update(name string, id int, data json.RawMessage) error
	Retrieve(name string, id int) (json.RawMessage, error)
	RetrieveAll(name string) (map[int]json.RawMessage, error)
}

// Environment provides the environment for this module.
type Environment interface {
	Host() string
	Port() string
	DBFilename() string
}
