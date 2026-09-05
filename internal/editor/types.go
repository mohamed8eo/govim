package editor

import (
	"github.com/mohamed8eo/govim/internal/buffer"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeCommand
	ModeOpending
)

type Editor struct {
	Buf       *buffer.Buffer
	Cx, Cy    int
	RowOff    int
	Mode      Mode
	CmdBuf    string
	Quit      bool
	StatusMsg string
}

type Commands struct {
	Save    bool
	Quit    bool
	Unknown bool
}
