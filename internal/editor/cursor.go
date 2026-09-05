package editor

func (ed *Editor) MoveCursor(dx, dy int) {
	ed.StatusMsg = ""
	ed.Cy += dy
	if ed.Cy < 0 {
		ed.Cy = 0
	}
	if ed.Cy >= ed.Buf.NumLine() {
		ed.Cy = ed.Buf.NumLine() - 1
	}

	ed.Cx += dx
	if ed.Cx < 0 {
		ed.Cx = 0
	}
	maxCx := max(0, len(ed.Buf.Line(ed.Cy)))
	if ed.Cx > maxCx {
		ed.Cx = maxCx
	}
}

func (ed *Editor) Scroll(textRows int) {
	if ed.Cy < ed.RowOff {
		ed.RowOff = ed.Cy
	}
	if ed.Cy >= ed.RowOff+textRows {
		ed.RowOff = ed.Cy - textRows + 1
	}
	if ed.RowOff < 0 {
		ed.RowOff = 0
	}
}
