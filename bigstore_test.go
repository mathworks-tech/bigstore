package bigstore

import (
	"os"
	"testing"
)

func TestBasicUsage(t *testing.T) {
	const testFile = "/tmp/bigstore_test.tick"
	os.Remove(testFile)
	defer os.Remove(testFile)

	store, err := New(BlockSz64)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Free()

	err = store.Open(testFile, ModeCreate, 1024*1024)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

}
