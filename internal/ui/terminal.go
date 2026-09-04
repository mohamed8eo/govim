package ui

import (
	"golang.org/x/sys/unix"
)

type Key int

const (
	KeyUp Key = iota + 1000
	KeyDown
	KeyLeft
	KeyRight
	KeyEsc
)

func ReadKey(fd int, raw, peek *unix.Termios) (rune, Key, error) {
	err := unix.IoctlSetTermios(fd, unix.TCSETS, raw)
	if err != nil {
		return 0, 0, err
	}

	buf := make([]byte, 1)
	n, err := unix.Read(fd, buf)
	if n <= 0 && err != nil {
		return 0, 0, err
	}

	if buf[0] == '\x1b' {
		err := unix.IoctlSetTermios(fd, unix.TCSETS, peek)
		if err != nil {
			return 0, 0, err
		}

		n, err = unix.Read(fd, buf)

		restErr := unix.IoctlSetTermios(fd, unix.TCSETS, raw)
		if restErr != nil {
			return 0, 0, restErr
		}

		if n <= 0 || err != nil {
			return 0, KeyEsc, nil
		}

		if buf[0] == '[' {
			n, err = unix.Read(fd, buf)
			if n <= 0 {
				return 0, 0, err
			}
			switch buf[0] {
			case 'A':
				return 0, KeyUp, nil
			case 'B':
				return 0, KeyDown, nil
			case 'C':
				return 0, KeyRight, nil
			case 'D':
				return 0, KeyLeft, nil
			}
		}
	}

	closeKey := buf[0]
	if buf[0] == 'q' || buf[0] == '\x03' {
		return rune(closeKey), 0, nil
	}

	return rune(buf[0]), 0, nil
}
