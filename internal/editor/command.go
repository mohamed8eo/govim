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

func (ed *Editor) ExecuteSearch() {
	ed.LastSearch = ed.SearchBuf
	ed.SearchBuf = ""
	ed.Mode = ModeNormal

	if ed.LastSearch == "" {
		return
	}

	ed.findNext()
}

func (ed *Editor) findNext() {
	if ed.LastSearch == "" {
		return
	}

	numLines := len(ed.Buf.Lines)
	startCy := ed.Cy

	// 1. rest of the current line, strictly after the cursor
	line := ed.Buf.Lines[ed.Cy]
	if ed.Cx+1 <= len(line) {
		rest := line[ed.Cx+1:]
		if idx := strings.Index(rest, ed.LastSearch); idx != -1 {
			ed.Cx = ed.Cx + 1 + idx
			return
		}
	}

	// 2. every line after this one, wrapping around to the top
	for offset := 1; offset <= numLines; offset++ {
		cy := (startCy + offset) % numLines
		line := ed.Buf.Lines[cy]
		if idx := strings.Index(line, ed.LastSearch); idx != -1 {
			ed.Cy = cy
			ed.Cx = idx
			return
		}
	}

	ed.StatusMsg = fmt.Sprintf("Pattern not found: %s", ed.LastSearch)
}

func (ed *Editor) CancelCommand() {
	switch ed.Mode {
	case ModeCommand:
		ed.CmdBuf = ""
	case ModeSearch:
		ed.SearchBuf = ""
	}
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
