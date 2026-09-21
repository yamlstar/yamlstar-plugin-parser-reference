// Package referenceparser exposes the generated YAML reference parser.
package referenceparser

import "github.com/yamlstar/yamlstar-plugin-parser-reference/parser"

// Event is one YAML parser event.
type Event = parser.Event

// Initialize loads the generated parser namespaces once.
func Initialize() error {
	return parser.Initialize()
}

// Parse returns YAML parser events without a serialization round trip.
func Parse(input []byte) ([]Event, error) {
	return parser.Parse(input)
}
