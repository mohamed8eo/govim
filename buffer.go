package main

import (
	"io"
	"os"
	"strings"
)

type Buffer struct {
	lines []string
	path  string
}

func LoadFile(path string) (*Buffer, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return &Buffer{path: path}, err
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return &Buffer{path: path}, err
	}

	data := string(buf)

	data = strings.TrimSuffix(data, "\n")
	dataSlice := strings.Split(data, "\n")

	return &Buffer{lines: dataSlice, path: path}, nil
}

func (b *Buffer) NumLine() int {
	return len(b.lines)
}

func (b *Buffer) Line(row int) string {
	return b.lines[row]
}
