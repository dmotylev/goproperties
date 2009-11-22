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
		t.Fatal("failed to open %s: ", _testFilename, err);
		return;
	}
	defer f.Close();
	Load(f);
}
