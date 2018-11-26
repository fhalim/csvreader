package csvreader

import (
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	reader, err := NewReader(strings.NewReader(`"key","value"
"answeroflife","42"
`))
	if err != nil {
		t.Error(err)
	}
	if len(reader.Columns) != 2 {
		t.Error("Invalid number of columns detected")
	}
	rowCount := 0
	for row := range reader.Data {
		rowCount += 1
		if row["key"] != "answeroflife" {
			t.Error("Could not find key column")
		}
		if row["value"] != "42" {
			t.Error("Could not find value column")
		}
	}
	if rowCount != 1 {
		t.Error("Invalid number of records", rowCount)
	}
}
