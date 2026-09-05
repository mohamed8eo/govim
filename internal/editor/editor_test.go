package editor

import (
	"testing"

	"github.com/mohamed8eo/govim/internal/buffer"
)

func TestNewEditor(t *testing.T) {
	buf := &buffer.Buffer{
		Lines: []string{"hello", "world"},
		Path:  "dummy.txt",
	}

	ed := NewEditor(buf)
	if ed.Mode != ModeNormal {
		t.Errorf("expected ModeNormal, got %v", ed.Mode)
	}
	if ed.Cx != 0 || ed.Cy != 0 {
		t.Errorf("expected cursor at (0,0), got (%d,%d)", ed.Cx, ed.Cy)
	}
}

func TestMoveCursor(t *testing.T) {
	buf := &buffer.Buffer{
		Lines: []string{"hello", "world"},
		Path:  "dummy.txt",
	}

	ed := NewEditor(buf)
	ed.MoveCursor(2, 1)

	if ed.Cy != 1 || ed.Cx != 2 {
		t.Errorf("expected cursor at (2,1), got (%d,%d)", ed.Cx, ed.Cy)
	}

	// Test bounds checking
	ed.MoveCursor(-10, -10)
	if ed.Cy != 0 || ed.Cx != 0 {
		t.Errorf("expected cursor at bounds (0,0), got (%d,%d)", ed.Cx, ed.Cy)
	}
}

func TestCheckCommand(t *testing.T) {
	cmd := checkCommand("wq")
	if !cmd.Save || !cmd.Quit || cmd.Unknown {
		t.Errorf("expected Save=true, Quit=true, Unknown=false, got Save=%v, Quit=%v, Unknown=%v", cmd.Save, cmd.Quit, cmd.Unknown)
	}

	cmdUnk := checkCommand("wxq")
	if !cmdUnk.Unknown {
		t.Errorf("expected Unknown=true for 'wxq'")
	}
}

func TestExecuteCommandUnknown(t *testing.T) {
	buf := &buffer.Buffer{
		Lines: []string{"hello"},
		Path:  "dummy.txt",
	}

	ed := NewEditor(buf)
	ed.Mode = ModeCommand
	ed.CmdBuf = "foo"
	ed.ExecuteCommand()

	if ed.Mode != ModeNormal {
		t.Errorf("expected ModeNormal after executing unknown command, got %v", ed.Mode)
	}
	if ed.StatusMsg != "unknown command: foo" {
		t.Errorf("expected 'unknown command: foo', got %q", ed.StatusMsg)
	}
}

func TestDeleteXEdgeCases(t *testing.T) {
	// Edge case 1: cursor at len(line) (insert append position after escape)
	buf1 := &buffer.Buffer{Lines: []string{"hello"}, Path: "dummy.txt"}
	ed1 := NewEditor(buf1)
	ed1.Cx = 5 // len("hello")
	ed1.DeleteX()
	if ed1.Buf.Lines[0] != "hell" {
		t.Errorf("expected line to be 'hell', got %q", ed1.Buf.Lines[0])
	}
	if ed1.Cx < 0 {
		t.Errorf("expected non-negative cursor, got %d", ed1.Cx)
	}

	// Edge case 2: single-char line "a", cx = 0
	buf2 := &buffer.Buffer{Lines: []string{"a"}, Path: "dummy.txt"}
	ed2 := NewEditor(buf2)
	ed2.Cx = 0
	ed2.DeleteX()
	if ed2.Buf.Lines[0] != "" {
		t.Errorf("expected line to be empty, got %q", ed2.Buf.Lines[0])
	}
	if ed2.Cx != 0 {
		t.Errorf("expected cursor at 0, got %d", ed2.Cx)
	}
}

func TestMoveToLineEndEmpty(t *testing.T) {
	buf := &buffer.Buffer{Lines: []string{""}, Path: "dummy.txt"}
	ed := NewEditor(buf)
	ed.MoveToLineEnd()
	if ed.Cx != 0 {
		t.Errorf("expected cx = 0 on empty line MoveToLineEnd, got %d", ed.Cx)
	}
}

func TestDeleteLineToEnd(t *testing.T) {
	buf := &buffer.Buffer{Lines: []string{"hello"}, Path: "dummy.txt"}
	ed := NewEditor(buf)
	ed.Cx = 2 // on 'l'
	ed.DeleteLineToEnd()
	if ed.Buf.Lines[0] != "he" {
		t.Errorf("expected line 'he', got %q", ed.Buf.Lines[0])
	}
	if ed.Cx != 1 {
		t.Errorf("expected cx = 1, got %d", ed.Cx)
	}
}
