package aibot

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type StreamDecoder struct {
	reader *bufio.Reader
	err    error
	done   bool
	field  string
	value  string
}

func NewDecoder(r io.Reader) *StreamDecoder {
	return &StreamDecoder{
		reader: bufio.NewReader(r),
	}
}

func (d *StreamDecoder) Next() bool {
	for !d.done && d.err == nil {
		var line []byte
		line, d.err = d.reader.ReadBytes('\n')
		if d.err != nil {
			return false
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == ':' {
			continue
		}

		parts := strings.SplitN(string(line), ":", 2)
		if len(parts) != 2 {
			continue
		}

		field := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if field == "" || value == "" {
			continue
		}

		switch field {
		case "data", "event":
			d.field = field
			d.value = value
			return true
		}
	}

	return false
}

func (d *StreamDecoder) Field() string {
	return d.field
}

func (d *StreamDecoder) Value() string {
	return d.value
}

func (d *StreamDecoder) Err() error {
	return d.err
}
