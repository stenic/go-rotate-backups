package drivers

import "testing"

func TestLocalDriverListDirsAllowsMissingPath(t *testing.T) {
	dirs, err := (&LocalDriver{}).ListDirs(t.TempDir() + "/missing")
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 0 {
		t.Fatalf("ListDirs() = %v, want empty list", dirs)
	}
}
