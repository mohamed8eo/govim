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

	for i := 0; i < textRows; i++ {
		fileRow := ed.RowOff + i
		if fileRow < len(ed.Buf.Lines) {
			line := ed.Buf.Lines[fileRow]
			if len(line) > cols {
				line = line[:cols]
			}
			sb.WriteString(line)
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
