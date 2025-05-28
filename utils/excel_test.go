package utils

import (
	"os"
	"testing"
)

func TestParseExcel(t *testing.T) {
	file, err := os.Open("../testdata/sample.xlsx")
	if err != nil {
		t.Fatalf("Failed to open sample Excel file: %v", err)
	}

	defer file.Close()

	rows, err := ParseExcel(file)
	if err != nil {
		t.Fatalf("ParseExcel returner error: %v", err)
	}

	if len(rows) < 2 {
		t.Fatalf("Expected at least 2 rows, got %d", len(rows))
	}

	if rows[1][0] != "Nurs" {
		t.Errorf("Expected 'Nurs' in first data row, got '%s'", rows[1][0])
	}
}
