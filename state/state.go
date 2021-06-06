package state

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	rw "github.com/mattn/go-runewidth"
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

type snapshot struct {
	text           []string
	anchor, cursor cursor
}

type diff struct {
	before, after              snapshot
	commonPrefix, commonSuffix int
}

type State struct {
	FilePath       string
	TabWidth       int
	Text           []string
	Anchor, Cursor cursor
	mode           mode
	Msg            string
	snapshot       snapshot
	history        []diff
	historyHead    int
}

func visualWidth(col, tabWidth int, char rune) int {
	if char == '\t' {
		return tabWidth - col%tabWidth
	}
	return rw.RuneWidth(char)
}

func (s *State) setCursorX(c *cursor, x int) {
	c.X = x
	c.col = 0
	for _, char := range s.Text[c.Y][:c.X] {
		c.col += visualWidth(c.col, s.TabWidth, char)
	}
}

func (s *State) setCursorY(c *cursor, y int) {
	c.Y = y
	c.X = -1
	col := 0
	var char rune
	for c.X, char = range s.Text[c.Y] {
		col += visualWidth(col, s.TabWidth, char)
		if col > c.col {
			break
		}
	}
}

func (s *State) moveLeft(c *cursor) {
	if c.X <= 0 {
		return
	}
	_, size := utf8.DecodeLastRuneInString(s.Text[c.Y][:c.X])
	s.setCursorX(c, c.X-size)
}

func (s *State) moveRight(c *cursor) {
	if len(s.Text[c.Y]) == 0 {
		return
	}
	_, sizeLast := utf8.DecodeLastRuneInString(s.Text[c.Y])
	if c.X == len(s.Text[c.Y])-sizeLast {
		return
	}
	_, size := utf8.DecodeRuneInString(s.Text[c.Y][c.X:])
	s.setCursorX(c, c.X+size)
}

func (s *State) moveUp(c *cursor, n int) {
	if c.Y == 0 {
		return
	}
	if c.Y-n >= 0 {
		s.setCursorY(c, c.Y-n)
	} else {
		s.setCursorY(c, 0)
	}
}

func (s *State) moveDown(c *cursor, n int) {
	if c.Y == len(s.Text)-1 {
		return
	}
	if c.Y+n <= len(s.Text)-1 {
		s.setCursorY(c, c.Y+n)
	} else {
		s.setCursorY(c, len(s.Text)-1)
	}
}

func (s *State) moveStartOfLine(c *cursor) {
	for x, char := range s.Text[c.Y] {
		if char != ' ' && char != '\t' {
			s.setCursorX(c, x)
			return
		}
	}
}

func (s *State) moveEndOfLine(c *cursor) {
	if len(s.Text[c.Y]) == 0 {
		return
	}
	_, sizeLast := utf8.DecodeLastRuneInString(s.Text[c.Y])
	s.setCursorX(c, len(s.Text[c.Y])-sizeLast)
}

var reStartOfWord = regexp.MustCompile(
	`([[:word:]]+|[[:punct:]]+)[[:blank:]]*$`,
)

func (s *State) moveStartOfWord(c *cursor) {
	if len(s.Text[c.Y]) == 0 {
		return
	}
	if match := reStartOfWord.FindStringIndex(
		s.Text[c.Y][:c.X],
	); match != nil {
		s.setCursorX(c, match[0])
	}
}

var reEndOfWord = regexp.MustCompile(
	`^[[:blank:]]*([[:word:]]+|[[:punct:]]+)`,
)

func (s *State) moveEndOfWord(c *cursor) {
	if len(s.Text[c.Y]) == 0 {
		return
	}
	if match := reEndOfWord.FindStringIndex(
		s.Text[c.Y][c.X+1:],
	); match != nil {
		s.setCursorX(c, c.X+match[1])
	}
}

func (s *State) move(f func(*cursor)) {
	f(&s.Cursor)
	s.Anchor = s.Cursor
}

func (s *State) newLineAbove() {
	s.Text = append(
		s.Text[:s.Cursor.Y],
		append([]string{""}, s.Text[s.Cursor.Y:]...)...,
	)
	s.setCursorY(&s.Cursor, s.Cursor.Y)
	s.Anchor = s.Cursor
}

func (s *State) normaliseSelection() {
	if s.Cursor.Y < s.Anchor.Y ||
		s.Cursor.Y == s.Anchor.Y && s.Cursor.X < s.Anchor.X {
		s.Cursor, s.Anchor = s.Anchor, s.Cursor
	}
}

func (s *State) delete() {
	s.normaliseSelection()
	if s.Anchor.X == -1 {
		s.Text[s.Cursor.Y] = s.Text[s.Cursor.Y][s.Cursor.X+1:]
	} else {
		s.Text[s.Cursor.Y] = s.Text[s.Anchor.Y][:s.Anchor.X] +
			s.Text[s.Cursor.Y][s.Cursor.X+1:]
	}
	s.Text = append(s.Text[:s.Anchor.Y], s.Text[s.Cursor.Y:]...)
	s.setCursorY(&s.Anchor, s.Anchor.Y)
	s.Cursor = s.Anchor
}

func (s *State) deleteLines() {
	s.normaliseSelection()
	s.Text = append(s.Text[:s.Anchor.Y], s.Text[s.Cursor.Y+1:]...)
	if s.Anchor.Y >= len(s.Text) {
		s.Anchor.Y = len(s.Text) - 1
	}
	s.setCursorY(&s.Anchor, s.Anchor.Y)
	s.Cursor = s.Anchor
}

func (s *State) insert(char rune) {
	if char == '\n' {
		newLine := s.Text[s.Cursor.Y][s.Cursor.X:]
		s.Text[s.Cursor.Y] = s.Text[s.Cursor.Y][:s.Cursor.X]
		s.Text = append(
			s.Text[:s.Cursor.Y+1],
			append([]string{newLine}, s.Text[s.Cursor.Y+1:]...)...,
		)
		s.setCursorY(&s.Cursor, s.Cursor.Y+1)
		s.setCursorX(&s.Cursor, 0)
		s.Anchor = s.Cursor
		return
	}
	s.Text[s.Cursor.Y] = s.Text[s.Cursor.Y][:s.Cursor.X] +
		string(char) + s.Text[s.Cursor.Y][s.Cursor.X:]
	s.setCursorX(&s.Cursor, s.Cursor.X+1)
	s.Anchor = s.Cursor
}

func (s *State) insertBackspace() {
	if s.Cursor.X == 0 && s.Cursor.Y == 0 {
		return
	} else if s.Cursor.X == 0 {
		x := len(s.Text[s.Cursor.Y-1])
		s.Text[s.Cursor.Y-1] += s.Text[s.Cursor.Y]
		s.Text = append(s.Text[:s.Cursor.Y], s.Text[s.Cursor.Y+1:]...)
		s.setCursorY(&s.Cursor, s.Cursor.Y-1)
		s.setCursorX(&s.Cursor, x)
		s.Anchor = s.Cursor
		return
	}
	s.Text[s.Cursor.Y] = s.Text[s.Cursor.Y][:s.Cursor.X-1] +
		s.Text[s.Cursor.Y][s.Cursor.X:]
	s.setCursorX(&s.Cursor, s.Cursor.X-1)
	s.Anchor = s.Cursor
}

func (s *State) insertBackspaceWord() {
	if s.Cursor.X == 0 && s.Cursor.Y == 0 {
		return
	} else if s.Cursor.X == 0 {
		s.insertBackspace()
		return
	} else if match := reStartOfWord.FindStringIndex(
		s.Text[s.Cursor.Y][:s.Cursor.X],
	); match != nil {
		s.Text[s.Cursor.Y] = s.Text[s.Cursor.Y][:match[0]] +
			s.Text[s.Cursor.Y][s.Cursor.X:]
		s.setCursorX(&s.Cursor, match[0])
		s.Anchor = s.Cursor
	} else {
		s.Text[s.Cursor.Y] = s.Text[s.Cursor.Y][s.Cursor.X:]
		s.setCursorX(&s.Cursor, 0)
		s.Anchor = s.Cursor
	}
}

func (s *State) insertDelete() {
	if s.Cursor.X >= len(s.Text[s.Cursor.Y]) && s.Cursor.Y == len(s.Text)-1 {
		return
	} else if s.Cursor.X >= len(s.Text[s.Cursor.Y]) {
		s.Text[s.Cursor.Y] += s.Text[s.Cursor.Y+1]
		s.Text = append(s.Text[:s.Cursor.Y+1], s.Text[s.Cursor.Y+2:]...)
		return
	}
	s.Text[s.Cursor.Y] = s.Text[s.Cursor.Y][:s.Cursor.X] +
		s.Text[s.Cursor.Y][s.Cursor.X+1:]
}

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

func (s *State) takeSnapshot() {
	s.snapshot.text = make([]string, len(s.Text))
	copy(s.snapshot.text, s.Text)
	s.snapshot.anchor = s.Anchor
	s.snapshot.cursor = s.Cursor
}

func (s *State) recordHistory() {
	var d diff

	for d.commonPrefix < len(s.Text) && d.commonPrefix < len(s.snapshot.text) &&
		s.Text[d.commonPrefix] == s.snapshot.text[d.commonPrefix] {
		d.commonPrefix++
	}

	if d.commonPrefix == len(s.Text) && d.commonPrefix == len(s.snapshot.text) {
		return
	}

	for d.commonSuffix < len(s.Text)-d.commonPrefix &&
		d.commonSuffix < len(s.snapshot.text)-d.commonPrefix &&
		s.Text[len(s.Text)-1-d.commonSuffix] ==
			s.snapshot.text[len(s.snapshot.text)-1-d.commonSuffix] {
		d.commonSuffix++
	}

	d.before.text = make(
		[]string,
		len(s.snapshot.text)-d.commonPrefix-d.commonSuffix,
	)
	copy(
		d.before.text,
		s.snapshot.text[d.commonPrefix:len(s.snapshot.text)-d.commonSuffix],
	)
	d.before.anchor = s.snapshot.anchor
	d.before.cursor = s.snapshot.cursor

	d.after.text = make(
		[]string,
		len(s.Text)-d.commonPrefix-d.commonSuffix,
	)
	copy(
		d.after.text,
		s.Text[d.commonPrefix:len(s.Text)-d.commonSuffix],
	)
	d.after.anchor = s.Anchor
	d.after.cursor = s.Cursor

	s.history = s.history[:s.historyHead]
	s.history = append(s.history, d)
	s.historyHead++
}

func (s *State) undo() {
	if s.historyHead == 0 {
		return
	}
	s.historyHead--
	h := s.history[s.historyHead]
	s.Text = append(
		s.Text[:h.commonPrefix],
		append(h.before.text, s.Text[len(s.Text)-h.commonSuffix:]...)...,
	)
	s.Anchor = h.before.anchor
	s.Cursor = h.before.cursor
}

func (s *State) redo() {
	if s.historyHead >= len(s.history) {
		return
	}
	h := s.history[s.historyHead]
	s.historyHead++
	s.Text = append(
		s.Text[:h.commonPrefix],
		append(h.after.text, s.Text[len(s.Text)-h.commonSuffix:]...)...,
	)
	s.Anchor = h.after.anchor
	s.Cursor = h.after.cursor
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
				s.takeSnapshot()
				s.setMode(modeInsert)
			case 'A':
				s.takeSnapshot()
				s.newLineAbove()
				s.setMode(modeInsert)
			case 'd':
				s.takeSnapshot()
				s.setCursorX(&s.Cursor, s.Cursor.X+1)
				s.setMode(modeInsert)
			case 'D':
				s.takeSnapshot()
				s.move(s.moveEndOfLine)
				s.setCursorX(&s.Cursor, s.Cursor.X+1)
				s.setMode(modeInsert)
				s.insert('\n')
			case 'f':
				s.takeSnapshot()
				s.delete()
				s.setMode(modeInsert)
			case 'F':
				s.takeSnapshot()
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
				s.takeSnapshot()
				s.delete()
				s.recordHistory()
			case 'X':
				s.takeSnapshot()
				s.deleteLines()
				s.recordHistory()
			case 'z':
				s.undo()
			case 'Z':
				s.redo()
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

// debug functions

func (s *State) debugUnicode() {
	if s.Cursor.X >= 0 {
		char, size := utf8.DecodeRuneInString(s.Text[s.Cursor.Y][s.Cursor.X:])
		s.Msg = fmt.Sprintf(
			"%#U %v %v",
			char, size, rw.RuneWidth(char),
		)
	} else {
		s.Msg = "EMPTY"
	}
}

func (s *State) debugCursor() {
	s.Msg = fmt.Sprintf("a:%+v c:%+v", s.Anchor, s.Cursor)
}
