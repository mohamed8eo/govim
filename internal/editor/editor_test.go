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
	if ed.statusMsg != "unknown command: foo" {
		t.Errorf("expected 'unknown command: foo', got %q", ed.statusMsg)
	}
}
