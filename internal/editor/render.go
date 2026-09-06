package editor

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

func (ed *Editor) Render() {
	var sb strings.Builder
	sb.WriteString("\x1b[2J\x1b[H")

	fd := int(os.Stdout.Fd())
	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	rows := 24
	cols := 80
	if err == nil && ws.Row > 0 {
		rows = int(ws.Row)
		cols = int(ws.Col)
	}

	textRows := rows - 1
	textRows = max(textRows, 1)

	ed.Scroll(textRows)

	var startY, startX, endY, endX int
	isVisual := ed.Mode == ModeVisual
	if isVisual {
		startY, startX, endY, endX = ed.SelectionRange()
	}

	for i := 0; i < textRows; i++ {
		fileRow := ed.RowOff + i
		if fileRow < len(ed.Buf.Lines) {
			line := ed.Buf.Lines[fileRow]
			if len(line) > cols {
				line = line[:cols]
			}
			if isVisual && fileRow >= startY && fileRow <= endY {
				if startY == endY {
					sX := startX
					eX := endX
					if sX > len(line) {
						sX = len(line)
					}
					if eX >= len(line) {
						eX = len(line) - 1
					}
					if sX < 0 {
						sX = 0
					}
					if eX < 0 {
						eX = -1
					}
					if sX <= eX && eX < len(line) {
						before := line[:sX]
						selected := line[sX : eX+1]
						after := line[eX+1:]
						sb.WriteString(before)
						sb.WriteString("\x1b[7m")
						sb.WriteString(selected)
						sb.WriteString("\x1b[0m")
						sb.WriteString(after)
					} else {
						sb.WriteString(line)
					}
				} else {
					if fileRow == startY {
						sX := startX
						if sX > len(line) {
							sX = len(line)
						}
						if sX < 0 {
							sX = 0
						}
						before := line[:sX]
						selected := line[sX:]
						sb.WriteString(before)
						sb.WriteString("\x1b[7m")
						sb.WriteString(selected)
						sb.WriteString("\x1b[0m")
					} else if fileRow > startY && fileRow < endY {
						sb.WriteString("\x1b[7m")
						sb.WriteString(line)
						sb.WriteString("\x1b[0m")
					} else if fileRow == endY {
						eX := endX
						if eX >= len(line) {
							eX = len(line) - 1
						}
						if eX < 0 {
							eX = -1
						}
						if eX >= 0 && eX < len(line) {
							selected := line[:eX+1]
							after := line[eX+1:]
							sb.WriteString("\x1b[7m")
							sb.WriteString(selected)
							sb.WriteString("\x1b[0m")
							sb.WriteString(after)
						} else {
							sb.WriteString("\x1b[7m")
							sb.WriteString(line)
							sb.WriteString("\x1b[0m")
						}
					}
				}
			} else {
				sb.WriteString(line)
			}
		} else {
			sb.WriteString("~")
		}
		sb.WriteString("\r\n")
	}

	var modeStr string
	switch ed.Mode {
	case ModeNormal:
		modeStr = "NORMAL "
	case ModeInsert:
		modeStr = "INSERT "
	case ModeCommand:
		modeStr = "COMMAND "
	case ModeOpending:
		modeStr = "O-PENDING"
	case ModeSearch:
		modeStr = "Searching"
	case ModeVisual:
		modeStr = "VISUAL"
	case ModeVisualLine:
		modeStr = "VISUAL LINE"
	}

	positionStr := fmt.Sprintf(" %d,%d ", ed.Cy+1, ed.Cx+1)
	filePath := ed.Buf.Path
	if filePath == "" {
		filePath = "[No Name]"
	}

	if ed.Mode == ModeCommand {
		cmdLine := fmt.Sprintf(":%s", ed.CmdBuf)
		if len(cmdLine) > cols {
			cmdLine = cmdLine[:cols]
		}
		sb.WriteString(cmdLine)
	} else if ed.Mode == ModeSearch {
		prefix := ed.SearchPrefix
		if prefix == 0 {
			prefix = '?'
		}
		searchLine := fmt.Sprintf("%c%s", prefix, ed.SearchBuf)
		if len(searchLine) > cols {
			searchLine = searchLine[:cols]
		}
		sb.WriteString(searchLine)
	} else if ed.StatusMsg != "" {
		msg := ed.StatusMsg
		if len(msg) > cols {
			msg = msg[:cols]
		}
		sb.WriteString(msg)
	} else {
		statusLeft := fmt.Sprintf("\x1b[7m %s \x1b[0m %s", modeStr, filePath)
		statusRight := fmt.Sprintf("\x1b[7m%s\x1b[0m", positionStr)

		leftLen := len(modeStr) + 2 + len(filePath) + 1
		rightLen := len(positionStr)
		padding := max(cols-leftLen-rightLen, 0)
		sb.WriteString(statusLeft)
		if padding > 0 {
			sb.WriteString(strings.Repeat(" ", padding))
		}
		sb.WriteString(statusRight)
	}

	sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", (ed.Cy-ed.RowOff)+1, ed.Cx+1))
	os.Stdout.WriteString(sb.String())
}

func (ed *Editor) SelectionRange() (startY, startX, endY, endX int) {
	sy, sx, ey, ex := ed.Vy, ed.Vx, ed.Cy, ed.Cx
	if sy > ey || (sy == ey && sx > ex) {
		sy, ey = ey, sy
		sx, ex = ex, sx
	}
	if ed.Mode == ModeVisualLine {
		eX := 0
		if ey >= 0 && ey < len(ed.Buf.Lines) {
			eX = len(ed.Buf.Lines[ey]) - 1
			if eX < 0 {
				eX = 0
			}
		}
		return sy, 0, ey, eX
	}
	if ed.Vy < ed.Cy || (ed.Vy == ed.Cy && ed.Vx <= ed.Cx) {
		return ed.Vy, ed.Vx, ed.Cy, ed.Cx
	}
	return ed.Cy, ed.Cx, ed.Vy, ed.Vx
}
