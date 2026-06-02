package utils

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureCSVFile_CreatesFileWhenMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.csv")
	header := []string{"id", "name", "email"}

	if err := EnsureCSVFile(path, header); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// File must exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected file to be created, but it does not exist")
	}

	// File must contain the header row
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open created file: %v", err)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 row (header), got %d", len(records))
	}
	if len(records[0]) != len(header) {
		t.Fatalf("header column count: got %d, want %d", len(records[0]), len(header))
	}
	for i, col := range header {
		if records[0][i] != col {
			t.Errorf("header[%d]: got %q, want %q", i, records[0][i], col)
		}
	}
}

func TestEnsureCSVFile_DoesNotOverwriteExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.csv")
	header := []string{"id", "name"}

	// Pre-create the file with extra data
	existing := "id,name\n1,Alice\n2,Bob\n"
	if err := os.WriteFile(path, []byte(existing), 0644); err != nil {
		t.Fatalf("failed to write existing file: %v", err)
	}

	// EnsureCSVFile must leave the file untouched
	if err := EnsureCSVFile(path, header); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(contents) != existing {
		t.Errorf("file was modified:\ngot:  %q\nwant: %q", string(contents), existing)
	}
}

func TestEnsureCSVFile_IdempotentOnRepeatCalls(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "idempotent.csv")
	header := []string{"a", "b", "c"}

	for i := 0; i < 3; i++ {
		if err := EnsureCSVFile(path, header); err != nil {
			t.Fatalf("call %d unexpected error: %v", i+1, err)
		}
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}
	// Still exactly one header row — not duplicated on repeat calls
	if len(records) != 1 {
		t.Errorf("expected 1 row after 3 calls, got %d", len(records))
	}
}

func TestEnsureCSVFile_DifferentHeaders(t *testing.T) {
	tests := []struct {
		name   string
		header []string
	}{
		{"user header", []string{"id", "name", "email", "password", "created_at"}},
		{"expense header", []string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}},
		{"single column", []string{"id"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "file.csv")

			if err := EnsureCSVFile(path, tt.header); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			f, err := os.Open(path)
			if err != nil {
				t.Fatalf("failed to open file: %v", err)
			}
			defer f.Close()

			records, _ := csv.NewReader(f).ReadAll()
			if len(records) != 1 {
				t.Fatalf("expected 1 header row, got %d", len(records))
			}
			if len(records[0]) != len(tt.header) {
				t.Errorf("column count: got %d, want %d", len(records[0]), len(tt.header))
			}
		})
	}
}
