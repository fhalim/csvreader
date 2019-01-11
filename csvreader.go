package csvreader

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log"
)

type Reader struct {
	Columns []string
	Data    <-chan map[string]string
}

func NewReader(ctx context.Context, reader io.Reader) (*Reader, error) {
	r := csv.NewReader(reader)
	header, err := r.Read()
	if err == io.EOF {
		return nil, errors.New("Unable to read file")
	}
	c := make(chan map[string]string)
	go func() {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				record, err := r.Read()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Println(err)
					return
				}
				c <- toMap(header, record)
			}

		}
	}()
	return &Reader{Columns: header, Data: c}, nil
}

func toMap(header []string, values []string) map[string]string {
	// TODO: Verify header and values are same length
	m := make(map[string]string, len(header))
	for i, col := range header {
		m[col] = values[i]
	}
	return m
}
