package editor

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mohamed8eo/govim/internal/buffer"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
)

type Editor struct {
	Buf    *buffer.Buffer
	Cx, Cy int
	Mode   Mode
	CmdBuf string
	Quit   bool
}

type Commands struct {
	Save    bool
	Quit    bool
	Unknown bool
}

func NewEditor(buf *buffer.Buffer) *Editor {
	return &Editor{
		Buf:  buf,
		Cx:   0,
		Cy:   0,
		Mode: ModeNormal,
	}
}

func (ed *Editor) Render() {
	var sb strings.Builder
	sb.WriteString("\x1b[2J\x1b[H")

	for _, line := range ed.Buf.Lines {
		sb.WriteString(line)
		sb.WriteString("\r\n")
	}

	var modeStr string
	switch ed.Mode {
	case ModeNormal:
		modeStr = "-- NORMAL --"
	case ModeInsert:
		modeStr = "-- INSERT --"
	case ModeCommand:
		modeStr = fmt.Sprintf(":%s", ed.CmdBuf)
	}
	sb.WriteString(modeStr)
	sb.WriteString("\r\n")

	sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", ed.Cy+1, ed.Cx+1))
	os.Stdout.WriteString(sb.String())
}

func (ed *Editor) MoveCursor(dx, dy int) {
	ed.Cy += dy
	if ed.Cy < 0 {
		ed.Cy = 0
	}
	if ed.Cy >= ed.Buf.NumLine() {
		ed.Cy = ed.Buf.NumLine() - 1
	}

	ed.Cx += dx
	if ed.Cx < 0 {
		ed.Cx = 0
	}
	maxCx := len(ed.Buf.Line(ed.Cy))
	if ed.Cx > maxCx {
		ed.Cx = maxCx
	}
}

func (ed *Editor) InsertChar(r rune) {
	if ed.Cy < 0 || ed.Cy >= len(ed.Buf.Lines) {
		return
	}
	line := ed.Buf.Lines[ed.Cy]
	if ed.Cx > len(line) {
		ed.Cx = len(line)
	}
	line = line[:ed.Cx] + string(r) + line[ed.Cx:]

	ed.Buf.Lines[ed.Cy] = line
	ed.Cx++
}

func (ed *Editor) InsertNewline() {
	if ed.Cy < 0 || ed.Cy >= len(ed.Buf.Lines) {
		return
	}
	line := ed.Buf.Lines[ed.Cy]
	if ed.Cx > len(line) {
		ed.Cx = len(line)
	}
	before := line[:ed.Cx]
	after := line[ed.Cx:]

	ed.Buf.Lines[ed.Cy] = before

	ed.Buf.Lines = append(ed.Buf.Lines[:ed.Cy+1], append([]string{after}, ed.Buf.Lines[ed.Cy+1:]...)...)

	ed.Cy++
	ed.Cx = 0
}

func (ed *Editor) DeleteBack() {
	if ed.Cx > 0 {
		if ed.Cy < 0 || ed.Cy >= len(ed.Buf.Lines) {
			return
		}
		line := ed.Buf.Lines[ed.Cy]
		if ed.Cx > len(line) {
			ed.Cx = len(line)
		}
		line = line[:ed.Cx-1] + line[ed.Cx:]
		ed.Buf.Lines[ed.Cy] = line

		ed.Cx--
		return
	}

	if ed.Cy <= 0 || ed.Cy >= len(ed.Buf.Lines) {
		return
	}

	current := ed.Buf.Lines[ed.Cy]
	prev := ed.Buf.Lines[ed.Cy-1]
	ed.Cx = len(prev)

	ed.Buf.Lines[ed.Cy-1] = prev + current

	ed.Buf.Lines = append(
		ed.Buf.Lines[:ed.Cy],
		ed.Buf.Lines[ed.Cy+1:]...,
	)

	ed.Cy--
}

func (ed *Editor) DeleteX() {
}

func (ed *Editor) ExecuteCommand() {
	cmd := checkCommand(ed.CmdBuf)
	if cmd.Unknown {
		ed.CmdBuf = ""
		ed.Mode = ModeNormal
		return
	}
	if cmd.Save {
		if err := ed.Save(); err != nil {
			log.Printf("Error saving file: %v\n", err)
		}
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

func (ed *Editor) Save() error {
	content := strings.Join(ed.Buf.Lines, "\n")
	err := os.WriteFile(
		ed.Buf.Path,
		[]byte(content),
		0o644,
	)
	if err != nil {
		return fmt.Errorf("failed to write file %q: %w", ed.Buf.Path, err)
	}
	return nil
}
