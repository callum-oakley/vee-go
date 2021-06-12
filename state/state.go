package state

import (
	"github.com/gdamore/tcell/v2"
)

type mode int

const (
	modeNormal mode = iota
	modeInsert
	modeSpace
)

type cursor struct {
	X, Y, col int
}

type State struct {
	FilePath       string
	TabWidth       int
	Text           []string
	Anchor, Cursor cursor
	mode           mode
	Msg            string
	change         change
	history        []change
	historyHead    int
}

func (s *State) HandleKey(e *tcell.EventKey) bool {
	switch s.mode {
	case modeNormal:
		switch e.Key() {
		case tcell.KeyRune:
			switch e.Rune() {
			// mode transitions
			case ' ':
				s.setMode(modeSpace)
			case 'a':
				s.startChange()
				s.setMode(modeInsert)
			case 'A':
				s.startChange()
				s.newLineAbove()
				s.setMode(modeInsert)
			case 'd':
				s.startChange()
				s.setCursorX(&s.Cursor, s.xRightOf(&s.Cursor))
				s.setMode(modeInsert)
			case 'D':
				s.startChange()
				s.move(s.moveEndOfLine)
				s.setCursorX(&s.Cursor, s.xRightOf(&s.Cursor))
				s.setMode(modeInsert)
				s.insert('\n')
			case 'f':
				s.startChange()
				s.delete()
				s.setMode(modeInsert)
			case 'F':
				s.startChange()
				s.deleteLines()
				s.newLineAbove()
				s.setMode(modeInsert)

			// movements
			case 'y':
				s.move(s.moveStartOfLine)
			case 'Y':
				s.moveStartOfLine(&s.Cursor)
			case 'o':
				s.move(s.moveEndOfLine)
			case 'O':
				s.moveEndOfLine(&s.Cursor)
			case 'u':
				s.move(s.moveStartOfWord)
			case 'U':
				s.moveStartOfWord(&s.Cursor)
			case 'i':
				s.move(s.moveEndOfWord)
			case 'I':
				s.moveEndOfWord(&s.Cursor)
			case 'h':
				s.move(s.moveLeft)
			case 'H':
				s.moveLeft(&s.Cursor)
			case 'l':
				s.move(s.moveRight)
			case 'L':
				s.moveRight(&s.Cursor)
			case 'k':
				s.move(func(c *cursor) { s.moveUp(c, 1) })
			case 'K':
				s.moveUp(&s.Cursor, 1)
			case 'j':
				s.move(func(c *cursor) { s.moveDown(c, 1) })
			case 'J':
				s.moveDown(&s.Cursor, 1)

			// actions
			case 'x':
				s.startChange()
				s.delete()
				s.endChange()
			case 'X':
				s.startChange()
				s.deleteLines()
				s.endChange()
			case 'z':
				s.undo()
			case 'Z':
				s.redo()
			case 'w':
				s.save()
			}
		case tcell.KeyUp:
			if e.Modifiers() == tcell.ModShift {
				s.moveUp(&s.Cursor, 9)
			} else {
				s.move(func(c *cursor) { s.moveUp(c, 9) })
			}
		case tcell.KeyDown:
			if e.Modifiers() == tcell.ModShift {
				s.moveDown(&s.Cursor, 9)
			} else {
				s.move(func(c *cursor) { s.moveDown(c, 9) })
			}
		case tcell.KeyESC:
			s.Anchor = s.Cursor
		}
	case modeInsert:
		switch e.Key() {
		case tcell.KeyRune:
			s.insert(e.Rune())
		case tcell.KeyTAB:
			s.insert('\t')
		case tcell.KeyCR:
			s.insert('\n')
		case tcell.KeyDEL:
			s.insertBackspace()
		case 0x17:
			s.insertBackspaceWord()
		case tcell.KeyDelete:
			s.insertDelete()
		case tcell.KeyESC:
			s.setMode(modeNormal)
		}
	case modeSpace:
		switch e.Key() {
		case tcell.KeyRune:
			switch e.Rune() {
			case 'q':
				return true
			}
		}
		s.mode = modeNormal
	}
	return false
}
