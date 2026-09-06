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

func (ed *Editor) maxNormalCx(line string) int {
	if len(line) == 0 {
		return 0
	}
	return len(line) - 1
}

func (ed *Editor) MoveToLineEnd() {
	if ed.Cy < 0 || ed.Cy >= len(ed.Buf.Lines) {
		return
	}
	line := ed.Buf.Lines[ed.Cy]
	if len(line) == 0 {
		ed.Cx = 0
		return
	}
	ed.Cx = ed.maxNormalCx(line)
}

// Yank
func (ed *Editor) YankLine() {
	line := ed.Buf.Lines[ed.Cy]
	ed.Register = line
	ed.RegisterLinewise = true
	ed.Mode = ModeNormal
}

func (ed *Editor) YankVisualLine() {
	startY, _, endY, _ := ed.SelectionRange()

	if startY > endY {
		startY, endY = endY, startY
	}

	var lines []string
	for y := startY; y <= endY; y++ {
		if y >= 0 && y < len(ed.Buf.Lines) {
			lines = append(lines, ed.Buf.Lines[y])
		}
	}

	ed.Register = strings.Join(lines, "\n")
	ed.RegisterLinewise = true
	ed.Mode = ModeNormal
}

func (ed *Editor) YankSelection() {
	startY, startX, endY, endX := ed.SelectionRange()

	if startY > endY || (startY == endY && startX > endX) {
		startY, endY = endY, startY
		startX, endX = endX, startX
	}

	// One line
	if startY == endY {
		line := ed.Buf.Lines[startY]

		ed.Register = line[startX:endX]
		ed.RegisterLinewise = false
		ed.Mode = ModeNormal
		return
	}

	// Multi-line
	var sb strings.Builder

	// First line
	sb.WriteString(ed.Buf.Lines[startY][startX:])
	sb.WriteByte('\n')

	// Middle lines
	for y := startY + 1; y < endY; y++ {
		sb.WriteString(ed.Buf.Lines[y])
		sb.WriteByte('\n')
	}

	// Last line
	sb.WriteString(ed.Buf.Lines[endY][:endX])

	ed.Register = sb.String()
	ed.RegisterLinewise = false
	ed.Mode = ModeNormal
}

func (ed *Editor) YankWord() {
	line := ed.Buf.Lines[ed.Cy]
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
		ed.Register = line[ed.Cx:]
	} else {
		ed.Register = line[ed.Cx : ed.Cx+idx]
	}
	ed.RegisterLinewise = false
	ed.Mode = ModeNormal
}

func (ed *Editor) Put() {
	if ed.Register == "" {
		return
	}

	line := ed.Buf.Lines[ed.Cy]

	// Linewise yank
	if ed.RegisterLinewise {
		linesToInsert := strings.Split(ed.Register, "\n")
		newLines := make([]string, 0, len(ed.Buf.Lines)+len(linesToInsert))
		newLines = append(newLines, ed.Buf.Lines[:ed.Cy+1]...)
		newLines = append(newLines, linesToInsert...)
		newLines = append(newLines, ed.Buf.Lines[ed.Cy+1:]...)
		ed.Buf.Lines = newLines

		ed.Cy += len(linesToInsert)
		ed.Cx = 0
		return
	}

	// Characterwise yank
	if ed.Cx > len(line) {
		ed.Cx = len(line)
	}

	line = line[:ed.Cx] + ed.Register + line[ed.Cx:]
	ed.Buf.Lines[ed.Cy] = line

	ed.Cx += len(ed.Register)
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
	if ed.Cy < 0 || ed.Cy >= len(ed.Buf.Lines) {
		return
	}
	line := ed.Buf.Lines[ed.Cy]
	if len(line) == 0 {
		return
	}

	if ed.Cx > ed.maxNormalCx(line) {
		ed.Cx = ed.maxNormalCx(line)
	}

	ed.Buf.Lines[ed.Cy] = line[:ed.Cx] + line[ed.Cx+1:]
	newLine := ed.Buf.Lines[ed.Cy]
	if ed.Cx > ed.maxNormalCx(newLine) {
		ed.Cx = ed.maxNormalCx(newLine)
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
	if ed.Cy < 0 || ed.Cy >= len(ed.Buf.Lines) {
		return
	}
	line := ed.Buf.Lines[ed.Cy]
	if len(line) == 0 {
		ed.Cx = 0
		return
	}
	if ed.Cx > ed.maxNormalCx(line) {
		ed.Cx = ed.maxNormalCx(line)
	}
	ed.Buf.Lines[ed.Cy] = line[:ed.Cx]
	newLine := ed.Buf.Lines[ed.Cy]
	ed.Cx = ed.maxNormalCx(newLine)
}
