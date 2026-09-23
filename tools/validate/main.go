// Command validate checks that every fixture under fixtures/v1/observation
// agrees with schemas/v1/observation.schema.json: fixtures under valid/ must
// pass, fixtures under invalid/ must fail. It is the CI gate that keeps the
// schema and its fixture suite honest as contributors change either one.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const (
	schemaPath    = "schemas/v1/observation.schema.json"
	fixturesRoot  = "fixtures/v1/observation"
	validSubdir   = "valid"
	invalidSubdir = "invalid"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "validate:", err)
		os.Exit(1)
	}
}

func run() error {
	return validateAll(schemaPath, fixturesRoot)
}

// validateAll is run's logic with schemaPath/fixturesRoot as parameters
// rather than the package constants, so tests can point it at the same
// files via paths relative to the test binary's own working directory
// (go test's cwd is the package directory, not the repo root `go run`
// uses).
func validateAll(schemaPath, fixturesRoot string) error {
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(schemaPath)
	if err != nil {
		return fmt.Errorf("compiling %s: %w", schemaPath, err)
	}

	var failures int
	failures += checkDir(schema, filepath.Join(fixturesRoot, validSubdir), true)
	failures += checkDir(schema, filepath.Join(fixturesRoot, invalidSubdir), false)

	if failures > 0 {
		return fmt.Errorf("%d fixture(s) did not validate as expected", failures)
	}
	return nil
}

// checkDir validates every *.json file in dir against schema, printing one
// report line per fixture, and returns the number of fixtures that did not
// match wantValid. It walks and opens files through an os.Root anchored at
// dir rather than plain filepath.WalkDir/os.Open, so a symlink inside a
// fixture directory can't be used to read a file outside it.
func checkDir(schema *jsonschema.Schema, dir string, wantValid bool) int {
	var failures int

	root, err := os.OpenRoot(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "walking %s: %v\n", dir, err)
		return failures + 1
	}
	defer func() { _ = root.Close() }()

	err = fs.WalkDir(root.FS(), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		reportPath := filepath.Join(dir, path)

		f, err := root.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL %s: open: %v\n", reportPath, err)
			failures++
			return nil
		}
		defer func() { _ = f.Close() }()

		instance, err := jsonschema.UnmarshalJSON(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FAIL %s: decode: %v\n", reportPath, err)
			failures++
			return nil
		}

		validateErr := schema.Validate(instance)
		got := validateErr == nil

		if got != wantValid {
			failures++
			if wantValid {
				fmt.Printf("FAIL %s: expected valid, got error: %v\n", reportPath, validateErr)
			} else {
				fmt.Printf("FAIL %s: expected invalid, but validated cleanly\n", reportPath)
			}
			return nil
		}

		fmt.Printf("ok   %s\n", reportPath)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "walking %s: %v\n", dir, err)
		failures++
	}
	return failures
}
