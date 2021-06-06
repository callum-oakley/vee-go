package state

import "fmt"

func setCursorShape(n int) {
	// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html
	// 1 blinking block (default)
	// 3 blinking underline
	fmt.Printf("\033[%d q", n)
}

func (s *State) setMode(m mode) {
	switch s.mode {
	case modeInsert:
		s.move(s.moveLeft)
		setCursorShape(1)
		s.recordHistory()
	}
	switch m {
	case modeInsert:
		if s.Cursor.X == -1 {
			s.setCursorX(&s.Cursor, 0)
		}
		s.Anchor = s.Cursor
		setCursorShape(3)
	}
	s.mode = m
}
