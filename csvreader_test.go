package csvreader

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	reader, err := NewReader(context.Background(), strings.NewReader(`"key","value"
"answeroflife","42"
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
	if rowCount != 2 {
		t.Error("Invalid number of records", rowCount)
	}
}
func TestCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	reader, err := NewReader(ctx, strings.NewReader(`"key","value"
"answeroflife","42"
"answeroflife","42"
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
		cancel()
		if row["key"] != "answeroflife" {
			t.Error("Could not find key column")
		}
		if row["value"] != "42" {
			t.Error("Could not find value column")
		}
	}
	if rowCount != 2 {
		t.Error("Invalid number of records. Should have bailed after the 2nd.", rowCount)
	}
}

type dummyReader struct {
	recordCount  int
	headerRead   bool
	currentIndex int
}

func (reader dummyReader) Read(p []byte) (n int, err error) {
	header := `"key","value"`
	record := `"answeroflife","42"`
	if !reader.headerRead {
		reader.headerRead = true
		return strings.NewReader(header).Read(p)
	}
	if reader.currentIndex < reader.recordCount {
		reader.currentIndex += 1
		return strings.NewReader(record).Read(p)
	} else {
		return 0, io.EOF
	}
}
func BenchmarkReads(b *testing.B) {
	ioReader := dummyReader{recordCount: b.N}
	reader, err := NewReader(context.Background(), ioReader)
	if err != nil {
		b.Error(err)
	}
	rowCount := 0
	for row := range reader.Data {
		rowCount += 1
		if row["key"] != "answeroflife" {
			b.Error("Could not find key column")
		}
		if row["value"] != "42" {
			b.Error("Could not find value column")
		}
	}
	if rowCount != 1 {
		b.Error("Invalid number of records", rowCount)
	}
}
