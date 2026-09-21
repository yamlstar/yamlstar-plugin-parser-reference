package parser_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/yamlstar/yamlstar-plugin-parser-reference/parser"
)

func TestParse(t *testing.T) {
	events, err := parser.Parse([]byte(
		"%YAML 1.2\n--- {a: &a [!!str 1, 'lambda'], b: *a}\n"))
	if err != nil {
		t.Fatal(err)
	}
	var scalars []string
	for _, event := range events {
		if event.Type == "scalar" {
			scalars = append(scalars, event.Value)
		}
	}
	if want := []string{"a", "1", "lambda", "b"}; !reflect.DeepEqual(scalars, want) {
		t.Fatalf("scalars: got %q, want %q", scalars, want)
	}
	if len(events) < 2 || events[1].Version != "1.2" {
		t.Fatalf("document version not preserved: %+v", events)
	}
}

func TestErrors(t *testing.T) {
	for _, input := range [][]byte{[]byte("["), {0xff}} {
		if _, err := parser.Parse(input); err == nil {
			t.Fatalf("accepted %q", input)
		}
	}
}

func TestConcurrent(t *testing.T) {
	const callers = 8
	var wait sync.WaitGroup
	errors := make(chan error, callers)
	for range callers {
		wait.Add(1)
		go func() {
			defer wait.Done()
			_, err := parser.Parse([]byte("a: true\n"))
			errors <- err
		}()
	}
	wait.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
}
