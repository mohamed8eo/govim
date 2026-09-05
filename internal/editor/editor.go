package editor

import (
	"github.com/mohamed8eo/govim/internal/buffer"
)

func NewEditor(buf *buffer.Buffer) *Editor {
	return &Editor{
		Buf:  buf,
		Cx:   0,
		Cy:   0,
		Mode: ModeNormal,
	}
}
