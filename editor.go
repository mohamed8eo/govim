package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
)

type Editor struct {
	buf    *Buffer
	cx, cy int
	mode   Mode
	cmdBuf string
	quit   bool
}
type Commands struct {
	save   bool
	quit   bool
	unkown bool
}

func NewEditor(buf *Buffer) *Editor {
	return &Editor{
		buf:  buf,
		cx:   0,
		cy:   0,
		mode: ModeNormal,
	}
}

func (ed *Editor) Render() {
	var sb strings.Builder
	sb.WriteString("\x1b[2J\x1b[H")

	for _, line := range ed.buf.lines {
		sb.WriteString(line)
		sb.WriteString("\r\n")
	}

	var modeStr string
	switch ed.mode {
	case ModeNormal:
		modeStr = "-- NORMAL --"
	case ModeInsert:
		modeStr = "-- INSERT --"
	case ModeCommand:
		modeStr = fmt.Sprintf(":%s", ed.cmdBuf)
	}
	sb.WriteString(modeStr)
	sb.WriteString("\r\n")

	sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", ed.cy+1, ed.cx+1))
	os.Stdout.WriteString(sb.String())
}

func (ed *Editor) MoveCursor(dx, dy int) {
	ed.cy += dy
	if ed.cy < 0 {
		ed.cy = 0
	}
	if ed.cy >= ed.buf.NumLine() {
		ed.cy = ed.buf.NumLine() - 1
	}

	ed.cx += dx
	if ed.cx < 0 {
		ed.cx = 0
	}
	maxCx := len(ed.buf.Line(ed.cy))
	if ed.cx > maxCx {
		ed.cx = maxCx
	}
}

func (ed *Editor) InsertChar(r rune) {
	line := ed.buf.lines[ed.cy]
	line = line[:ed.cx] + string(r) + line[ed.cx:]

	ed.buf.lines[ed.cy] = line
	ed.cx++
}

func (ed *Editor) InsertNewline() {
	line := ed.buf.lines[ed.cy]
	before := line[:ed.cx]
	after := line[ed.cx:]

	ed.buf.lines[ed.cy] = before

	ed.buf.lines = append(ed.buf.lines[:ed.cy+1], append([]string{after}, ed.buf.lines[ed.cy+1:]...)...)

	ed.cy++
	ed.cx = 0
}

func (ed *Editor) DeleteBack() {
	if ed.cx > 0 {
		line := ed.buf.lines[ed.cy]
		line = line[:ed.cx-1] + line[ed.cx:]
		ed.buf.lines[ed.cy] = line

		ed.cx--
		return
	}

	if ed.cy == 0 {
		return
	}

	current := ed.buf.lines[ed.cy]
	prev := ed.buf.lines[ed.cy-1]
	ed.cx = len(prev)

	// join two line
	ed.buf.lines[ed.cy-1] = prev + current

	ed.buf.lines = append(
		ed.buf.lines[:ed.cy],
		ed.buf.lines[ed.cy+1:]...,
	)

	ed.cy--
}

func (ed *Editor) DeleteX() {
}

func (ed *Editor) ExecuteCommand() {
	cmd := checkCommand(ed.cmdBuf)
	if cmd.unkown {
		ed.cmdBuf = ""
		ed.mode = ModeNormal
		return
	}
	if cmd.save {
		err := ed.Save()
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
		}

	}
	if cmd.quit {
		ed.quit = true
		return
	}

	ed.cmdBuf = ""
	ed.mode = ModeNormal
}

func (ed *Editor) CancelCommand() {
	ed.cmdBuf = ""
	ed.mode = ModeNormal
}

func checkCommand(cmdBuf string) *Commands {
	var cmd Commands
	for _, command := range cmdBuf {
		switch command {
		case 'w':
			cmd.save = true
		case 'q':
			cmd.quit = true

		default:
			cmd.unkown = true
		}
	}

	return &cmd
}

func (ed *Editor) Save() error {
	content := strings.Join(ed.buf.lines, "\n")

	err := os.WriteFile(
		ed.buf.path,
		[]byte(content),
		0o644,
	)

	return err
}
