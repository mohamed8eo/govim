package main

import (
	"errors"
	"io"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("File path is required")
	}
	buffer, err := LoadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	ed := NewEditor(buffer)
	ed.Render()

	fd := int(os.Stdin.Fd())

	origTermios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		panic(err)
	}

	raw := *origTermios
	//&^ <- andNot gate
	raw.Lflag &^= unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN
	raw.Iflag &^= unix.IXON | unix.ICRNL
	raw.Oflag &^= unix.OPOST
	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0

	peek := raw
	peek.Cc[unix.VMIN] = 0
	peek.Cc[unix.VTIME] = 1

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &raw); err != nil {
		panic(err)
	}
	defer unix.IoctlSetTermios(fd, unix.TCSETS, origTermios)
	for {
		pressKey, keyModel, err := readKey(fd, &raw, &peek)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Printf("Error: %s", err.Error())
			continue
		}

		switch ed.mode {
		case ModeNormal:
			switch {
			case keyModel == KeyLeft, pressKey == 'h':
				ed.MoveCursor(-1, 0)
			case keyModel == KeyRight, pressKey == 'l':
				ed.MoveCursor(1, 0)
			case keyModel == KeyDown, pressKey == 'j':
				ed.MoveCursor(0, 1)
			case keyModel == KeyUp, pressKey == 'k':
				ed.MoveCursor(0, -1)
			case pressKey == 'a':
				ed.MoveCursor(1, 0)
				ed.mode = ModeInsert
			case pressKey == 'i':
				ed.mode = ModeInsert
			case pressKey == 'x':
				ed.DeleteX()
			case pressKey == '\x13':
				err := ed.Save()
				if err != nil {
					log.Fatalf("Error: %s\n", err.Error())
				}
			case pressKey == ':':
				ed.mode = ModeCommand
				ed.cmdBuf = ""
			}
		case ModeInsert:
			switch {
			case keyModel == KeyEsc:
				ed.mode = ModeNormal
			case pressKey == '\r':
				ed.InsertNewline()
			case pressKey == 127 || pressKey == '\b':
				ed.DeleteBack()
			case pressKey != 0:
				ed.InsertChar(pressKey)
			}
		case ModeCommand:
			switch {
			case pressKey == '\r':
				ed.ExecuteCommand()
			case keyModel == KeyEsc:
				ed.CancelCommand()

			default:
				ed.cmdBuf += string(pressKey)

			}
		}
		ed.Render()
		if ed.quit {
			return
		}

	}
}
