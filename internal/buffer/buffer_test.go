package buffer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "govim-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	buf, err := LoadFile(tmpFile)
	if err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}

	if buf.NumLine() != 3 {
		t.Errorf("expected 3 lines, got %d", buf.NumLine())
	}

	if buf.Line(0) != "line1" {
		t.Errorf("expected line 0 to be 'line1', got %q", buf.Line(0))
	}

	if buf.Line(2) != "line3" {
		t.Errorf("expected line 2 to be 'line3', got %q", buf.Line(2))
	}
}
