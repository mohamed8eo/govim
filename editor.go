package main

import (
	"fmt"
	"os"
	"strings"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
)

type Editor struct {
	buf    *Buffer
	cx, cy int
	mode   Mode
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
	modeStr := "-- NORMAL --"
	if ed.mode == ModeInsert {
		modeStr = "-- INSERT --"
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
