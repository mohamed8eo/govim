package editor

import (
	"strings"
)

func (ed *Editor) InsertChar(r rune) {
	ed.StatusMsg = ""
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
	ed.StatusMsg = ""
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

func (ed *Editor) MoveToLineEnd() {
	line := ed.Buf.Lines[ed.Cy]

	if ed.Cx <= len(line) {
		ed.Cx = len(line) - 1
	}
}

// Delete
func (ed *Editor) DeleteBack() {
	ed.StatusMsg = ""
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
	ed.StatusMsg = ""
	line := ed.Buf.Lines[ed.Cy]

	if len(line) == 0 {
		return
	}

	if ed.Cx > len(line) {
		ed.Cx = len(line) - 1
	}

	ed.Buf.Lines[ed.Cy] = line[:ed.Cx] + line[ed.Cx+1:]
	if ed.Cx == len(line)-1 {
		ed.Cx--
	}
}

func (ed *Editor) DeleteWord() {
	line := ed.Buf.Lines[ed.Cy]
	if ed.Cx >= len(line) {
		ed.Mode = ModeNormal
		return
	}
	wordStr := line[ed.Cx:]

	idx := strings.IndexFunc(wordStr, func(r rune) bool {
		switch r {
		case ' ', '\n', '\t', '-', '.', ',':
			return true
		default:
			return false
		}
	})

	if idx == -1 {
		ed.Buf.Lines[ed.Cy] = line[:ed.Cx]
	} else {
		ed.Buf.Lines[ed.Cy] = line[:ed.Cx] + line[ed.Cx+idx:]
	}

	ed.Mode = ModeNormal
}

func (ed *Editor) DeleteLine() {
	ed.StatusMsg = ""
	if len(ed.Buf.Lines) == 0 {
		return
	}

	up := ed.Buf.Lines[:ed.Cy]
	down := ed.Buf.Lines[ed.Cy+1:]
	ed.Buf.Lines = append(up, down...)

	if ed.Cy >= len(ed.Buf.Lines) {
		ed.Cy = len(ed.Buf.Lines) - 1
	}
	if ed.Cy < 0 {
		ed.Cy = 0
	}

	// Clamp horizontal cursor position to new line length
	maxCx := len(ed.Buf.Line(ed.Cy)) - 1
	maxCx = max(maxCx, 0)
	if ed.Cx > maxCx {
		ed.Cx = maxCx
	}

	ed.Mode = ModeNormal
}

func (ed *Editor) DeleteLineToEnd() {
	ed.CmdBuf = ""
	if len(ed.Buf.Lines) == 0 {
		return
	}

	line := ed.Buf.Lines[ed.Cy]
	ed.Buf.Lines[ed.Cy] = line[:ed.Cx]

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
}
