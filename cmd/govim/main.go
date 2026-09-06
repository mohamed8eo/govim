package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mohamed8eo/govim/internal/buffer"
	"github.com/mohamed8eo/govim/internal/editor"
	"github.com/mohamed8eo/govim/internal/ui"
	"golang.org/x/sys/unix"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Fprintln(os.Stderr, "Error: File path argument is required.")
		fmt.Fprintln(os.Stderr, "Usage: govim <filename>")
		os.Exit(1)
	}
	filePath := os.Args[1]

	buf, err := buffer.LoadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading file: %v\n", err)
		os.Exit(1)
	}

	// Switch to alternate screen buffer
	os.Stdout.WriteString("\x1b[?1049h")
	defer func() {
		// restore alternate screen
		os.Stdout.WriteString("\x1b[?1049l")
	}()

	ed := editor.NewEditor(buf)
	ed.Render()

	fd := int(os.Stdin.Fd())

	origTermios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting terminal attributes: %v\n", err)
		os.Exit(1)
	}

	raw := *origTermios
	raw.Lflag &^= unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN
	raw.Iflag &^= unix.IXON | unix.ICRNL
	raw.Oflag &^= unix.OPOST
	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0

	peek := raw
	peek.Cc[unix.VMIN] = 0
	peek.Cc[unix.VTIME] = 1

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &raw); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting raw terminal mode: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if rErr := unix.IoctlSetTermios(fd, unix.TCSETS, origTermios); rErr != nil {
			fmt.Fprintf(os.Stderr, "Error restoring terminal attributes: %v\n", rErr)
		}
	}()

	for {
		pressKey, keyModel, err := ui.ReadKey(fd, &raw, &peek, ed.Mode)
		if err != nil {
			log.Printf("Warning: error reading key: %v", err)
			continue
		}

		switch ed.Mode {
		case editor.ModeNormal:
			switch {
			case keyModel == ui.KeyLeft, pressKey == 'h':
				ed.MoveCursor(-1, 0)
			case keyModel == ui.KeyRight, pressKey == 'l':
				ed.MoveCursor(1, 0)
			case keyModel == ui.KeyDown, pressKey == 'j':
				ed.MoveCursor(0, 1)
			case keyModel == ui.KeyUp, pressKey == 'k':
				ed.MoveCursor(0, -1)
			case pressKey == 'a':
				ed.MoveCursor(1, 0)
				ed.Mode = editor.ModeInsert
			case pressKey == 'i':
				ed.Mode = editor.ModeInsert
			case pressKey == 'd':
				ed.PendingOp = 'd'
				ed.Mode = editor.ModeOpending
			case pressKey == 'y':
				ed.PendingOp = 'y'
				ed.Mode = editor.ModeOpending
			case pressKey == 'p':
				ed.Put()
			case pressKey == ':':
				ed.Mode = editor.ModeCommand
				ed.CmdBuf = ""
			case pressKey == '/', pressKey == '?':
				ed.Mode = editor.ModeSearch
				ed.SearchPrefix = pressKey
				ed.SearchBuf = ""
			case pressKey == 'v':
				ed.Mode = editor.ModeVisual
				ed.Vx, ed.Vy = ed.Cx, ed.Cy
			case pressKey == 'V':
				ed.Mode = editor.ModeVisualLine
				ed.Vx, ed.Vy = ed.Cx, ed.Cy
			case pressKey == 'x':
				ed.DeleteX()
			case pressKey == '\x13': // ctrl s
				ed.SaveWithStatus()
			case pressKey == 'D':
				ed.DeleteLineToEnd()
			case pressKey == '$':
				ed.MoveToLineEnd()
			case pressKey == '0':
				ed.Cx = 0

			}
		case editor.ModeInsert:
			switch {
			case keyModel == ui.KeyEsc:
				ed.Mode = editor.ModeNormal
			case pressKey == '\r':
				ed.InsertNewline()
			case pressKey == 127 || pressKey == '\b':
				ed.DeleteBack()
			case pressKey != 0:
				ed.InsertChar(pressKey)
			}
		case editor.ModeCommand:
			switch {
			case pressKey == '\r':
				ed.ExecuteCommand()
			case keyModel == ui.KeyEsc:
				ed.CancelCommand()
			default:
				ed.CmdBuf += string(pressKey)
			}

		case editor.ModeOpending:
			switch {
			case keyModel == ui.KeyEsc:
				ed.Mode = editor.ModeNormal
			case ed.PendingOp == 'd' && pressKey == 'd':
				ed.DeleteLine()
			case ed.PendingOp == 'd' && pressKey == 'w':
				ed.DeleteWord()
			case ed.PendingOp == 'y' && pressKey == 'y':
				ed.YankLine()
			case ed.PendingOp == 'y' && pressKey == 'w':
				ed.YankWord()
			default:
				ed.Mode = editor.ModeNormal
			}

		case editor.ModeSearch:
			switch {
			case pressKey == '\r':
				ed.ExecuteSearch()
			case keyModel == ui.KeyEsc:
				ed.CancelCommand()
			default:
				ed.SearchBuf += string(pressKey)
			}

		case editor.ModeVisual:
			switch {
			case keyModel == ui.KeyLeft, pressKey == 'h':
				ed.MoveCursor(-1, 0)
			case keyModel == ui.KeyRight, pressKey == 'l':
				ed.MoveCursor(1, 0)
			case keyModel == ui.KeyDown, pressKey == 'j':
				ed.MoveCursor(0, 1)
			case keyModel == ui.KeyUp, pressKey == 'k':
				ed.MoveCursor(0, -1)
			case pressKey == 'y':
				ed.YankSelection()

			case keyModel == ui.KeyEsc:
				ed.Mode = editor.ModeNormal
			}

		case editor.ModeVisualLine:
			switch {
			case keyModel == ui.KeyDown, pressKey == 'j':
				ed.MoveCursor(0, 1)
			case keyModel == ui.KeyUp, pressKey == 'k':
				ed.MoveCursor(0, -1)
			case pressKey == 'y':
				ed.YankVisualLine()

			case keyModel == ui.KeyEsc:
				ed.Mode = editor.ModeNormal
			}

		}
		ed.Render()
		if ed.Quit {
			return
		}
	}
}
