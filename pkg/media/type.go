package media

import (
	"strings"
)

// Type represents a general media content type.
type Type string

// String returns the type as string.
func (t Type) String() string {
	return string(t)
}

// Equal checks if the type matches.
func (t Type) Equal(s string) bool {
	return strings.EqualFold(s, t.String())
}

// NotEqual checks if the type is different.
func (t Type) NotEqual(s string) bool {
	return !t.Equal(s)
}

// Unknown checks if the type is unknown.
func (t Type) Unknown() bool {
	return t == Unknown
}
