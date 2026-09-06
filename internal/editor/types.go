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
	ModeSearch
	ModeVisual
	ModeVisualLine
)

type Editor struct {
	Buf                   *buffer.Buffer
	Cx, Cy                int
	Vx, Vy                int
	RowOff                int
	Mode                  Mode
	CmdBuf                string
	Quit                  bool
	StatusMsg             string
	SearchBuf, LastSearch string
	SearchPrefix          rune
	PendingOp             rune
	Register              string
	RegisterLinewise      bool
}

type Commands struct {
	Save    bool
	Quit    bool
	Unknown bool
}
