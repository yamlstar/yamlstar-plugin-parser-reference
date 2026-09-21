// Package parser exposes the generated YAML reference parser.
package parser

import (
	"fmt"
	"sync"
	"unicode/utf8"

	"github.com/glojurelang/glojure/pkg/glj"
	"github.com/glojurelang/glojure/pkg/lang"

	_ "github.com/yamlstar/yamlstar-plugin-parser-reference/internal/glojure/pkg/yaml_parser/core"
	_ "github.com/yamlstar/yamlstar-plugin-parser-reference/internal/glojure/pkg/yaml_parser/grammar"
	_ "github.com/yamlstar/yamlstar-plugin-parser-reference/internal/glojure/pkg/yaml_parser/parser"
	_ "github.com/yamlstar/yamlstar-plugin-parser-reference/internal/glojure/pkg/yaml_parser/prelude"
	_ "github.com/yamlstar/yamlstar-plugin-parser-reference/internal/glojure/pkg/yaml_parser/receiver"
)

// Version is the reference parser release version.
const Version = "0.2.5"

var (
	parseMu        sync.Mutex
	initializeOnce sync.Once
	initializeErr  error
)

// Event describes a YAMLStar parser event using ordinary Go values.
// Style is empty for plain scalars. Name is used only by aliases.
// The reference parser does not supply source positions or comments.
type Event struct {
	Type, Value, Anchor, Tag, Style, Name, Version string
	Flow, Explicit                                 bool
}

// Initialize loads the generated namespaces once.
func Initialize() error {
	initializeOnce.Do(func() {
		defer func() {
			if value := recover(); value != nil {
				initializeErr = fmt.Errorf("initialize plugin: %v", value)
			}
		}()
		require := glj.Var("clojure.core", "require")
		for _, namespace := range []string{
			"yaml-parser.prelude",
			"yaml-parser.grammar",
			"yaml-parser.parser",
			"yaml-parser.receiver",
			"yaml-parser.core",
		} {
			require.Invoke(lang.NewSymbol(namespace))
		}
	})
	return initializeErr
}

// Parse returns reference-parser events without a serialization round trip.
// Concurrent calls are supported, but parsing is serialized because the
// generated parser contains mutable namespace state.
func Parse(input []byte) (events []Event, err error) {
	parseMu.Lock()
	defer parseMu.Unlock()
	if !utf8.Valid(input) {
		return nil, fmt.Errorf("parser input must be UTF-8")
	}
	if err := Initialize(); err != nil {
		return nil, err
	}
	defer func() {
		if value := recover(); value != nil {
			events = nil
			if cause, ok := value.(error); ok {
				err = fmt.Errorf("parse: %w", cause)
			} else {
				err = fmt.Errorf("parse: %v", value)
			}
		}
	}()
	result := glj.Var("yaml-parser.core", "parse").Invoke(string(input))
	for seq := lang.Seq(result); seq != nil; seq = seq.Next() {
		item := seq.First()
		field := func(key string) string {
			value := lang.Get(item, lang.NewKeyword(key))
			if value == nil {
				return ""
			}
			text, ok := value.(string)
			if !ok {
				panic(fmt.Errorf("unexpected %s field type %T", key, value))
			}
			return text
		}
		flag := func(key string) bool {
			value := lang.Get(item, lang.NewKeyword(key))
			if value == nil {
				return false
			}
			flag, ok := value.(bool)
			if !ok {
				panic(fmt.Errorf("unexpected %s field type %T", key, value))
			}
			return flag
		}
		events = append(events, Event{
			Type: field("event"), Value: field("value"),
			Anchor: field("anchor"), Tag: field("tag"),
			Style: field("style"), Name: field("name"),
			Version: field("version"), Flow: flag("flow"),
			Explicit: flag("explicit"),
		})
	}
	return events, nil
}
