// test program for the props package
package properties

import (
	"testing";
	"os";
)

const _testFilename = "test.properties"

func TestLoad(t *testing.T) {
	f, err := os.Open(_testFilename, os.O_RDONLY, 0);
	if err != nil {
		t.Fatalf("failed to open %s: %s", _testFilename, err);
		return;
	}
	defer f.Close();
	props, loadErr := Load(f);
	if loadErr != nil {
		t.Fatalf("failed to load %s: %s", _testFilename, loadErr)
	}

	if props["website"] != "http://en.wikipedia.org/" {
		t.Error("website")
	}
	if props["language"] != "English" {
		t.Error("language")
	}
	if props["message"] != "Welcome to Wikipedia!" {
		t.Error("message")
	}
	if props["unicode"] != "\u041f\u0440\u0438\u0432\u0435\u0442, \u0421\u043e\u0432\u0430!" {
		t.Error("unicode")
	}
	if props["key with spaces"] != "This is the value that could be looked up with the key \"key with spaces\"." {
		t.Error("key with spaces")
	}
}
