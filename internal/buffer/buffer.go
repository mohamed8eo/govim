package buffer

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Buffer struct {
	Lines []string
	Path  string
}

func LoadFile(path string) (*Buffer, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return &Buffer{Path: path}, fmt.Errorf("failed to open file %q: %w", path, err)
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return &Buffer{Path: path}, fmt.Errorf("failed to read file %q: %w", path, err)
	}

	data := string(buf)
	data = strings.TrimSuffix(data, "\n")
	dataSlice := strings.Split(data, "\n")

	return &Buffer{Lines: dataSlice, Path: path}, nil
}

func (b *Buffer) NumLine() int {
	return len(b.Lines)
}

func (b *Buffer) Line(row int) string {
	if row < 0 || row >= len(b.Lines) {
		return ""
	}
	return b.Lines[row]
}
