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
	//for i := 0; i < len(props); i++ {
	t.Logf("%v\n", props);
	//}
}
