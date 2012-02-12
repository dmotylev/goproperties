// test program for the properties package
package goproperties

import (
	"os"
	"testing"
)

const _testFilename = "test.properties"

func TestLoad(t *testing.T) {
	f, err := os.Open(_testFilename)
	if err != nil {
		t.Fatalf("failed to open %s: %s", _testFilename, err)
		return
	}
	defer f.Close()
	props, loadErr := Load(f)
	if loadErr != nil {
		t.Fatalf("failed to load %s: %s", _testFilename, loadErr)
	}

	testValue(t, "website", "http://en.wikipedia.org/", props["website"])
	testValue(t, "language", "English", props["language"])
	testValue(t, "message", "Welcome to Wikipedia!", props["message"])
	testValue(t, "unicode", "Привет, Сова!", props["unicode"])
	testValue(t, "key with spaces", "This is the value that could be looked up with the key \"key with spaces\".", props["key with spaces"])
}

func testValue(t *testing.T, key, expected, value string) {
	if value != expected {
		t.Errorf("key     : '%s'\nexpected: '%s'\nvalue   : '%s'\n", key, expected, value)
	}
}
