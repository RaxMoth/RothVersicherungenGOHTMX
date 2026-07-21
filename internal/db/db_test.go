package db

import (
	"path/filepath"
	"testing"
)

func TestMigrate(t *testing.T) {
	database, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	// Running migrations twice must be a no-op.
	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate (second run): %v", err)
	}
}
