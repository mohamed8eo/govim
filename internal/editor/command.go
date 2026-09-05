package editor

import (
	"fmt"
	"os"
	"strings"
)

func (ed *Editor) ExecuteCommand() {
	cmd := checkCommand(ed.CmdBuf)
	if cmd.Unknown {
		ed.StatusMsg = fmt.Sprintf("unknown command: %s", ed.CmdBuf)
		ed.CmdBuf = ""
		ed.Mode = ModeNormal
		return
	}
	if cmd.Save {
		ed.SaveWithStatus()
	}
	if cmd.Quit {
		ed.Quit = true
		return
	}

	ed.CmdBuf = ""
	ed.Mode = ModeNormal
}

func (ed *Editor) CancelCommand() {
	ed.CmdBuf = ""
	ed.Mode = ModeNormal
}

func (ed *Editor) SaveWithStatus() {
	if err := save(ed.Buf.Lines, ed.Buf.Path); err != nil {
		ed.StatusMsg = fmt.Sprintf("Error saving file: %v", err)
	} else {
		ed.StatusMsg = fmt.Sprintf("%q written", ed.Buf.Path)
	}
}

func checkCommand(cmdBuf string) *Commands {
	var cmd Commands
	for _, command := range cmdBuf {
		switch command {
		case 'w':
			cmd.Save = true
		case 'q':
			cmd.Quit = true
		default:
			cmd.Unknown = true
		}
	}
	return &cmd
}

func save(lines []string, path string) error {
	content := strings.Join(lines, "\n")
	err := os.WriteFile(
		path,
		[]byte(content),
		0o644,
	)
	if err != nil {
		return fmt.Errorf("failed to write file %q: %w", path, err)
	}
	return nil
}
