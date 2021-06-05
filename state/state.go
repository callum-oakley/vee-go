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
	modeSpace
)

type cursor struct {
	X, Y, col int
}

type State struct {
	FilePath       string
	TabWidth       int
	Text           []string
	Cursor, Anchor cursor
	mode           mode
	Msg            string
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

func (s *State) HandleKey(e *tcell.EventKey) bool {
	switch s.mode {
	case modeNormal:
		switch e.Key() {
		case tcell.KeyRune:
			switch e.Rune() {
			case ' ':
				s.mode = modeSpace
			case 'y':
				s.move(s.moveStartOfLine)
			case 'o':
				s.move(s.moveEndOfLine)
			case 'u':
				s.move(s.moveStartOfWord)
			case 'i':
				s.move(s.moveEndOfWord)
			case 'h':
				s.move(s.moveLeft)
			case 'l':
				s.move(s.moveRight)
			case 'k':
				s.move(func(c *cursor) { s.moveUp(c, 1) })
			case 'j':
				s.move(func(c *cursor) { s.moveDown(c, 1) })
			case 'Y':
				s.moveStartOfLine(&s.Cursor)
			case 'O':
				s.moveEndOfLine(&s.Cursor)
			case 'U':
				s.moveStartOfWord(&s.Cursor)
			case 'I':
				s.moveEndOfWord(&s.Cursor)
			case 'H':
				s.moveLeft(&s.Cursor)
			case 'L':
				s.moveRight(&s.Cursor)
			case 'K':
				s.moveUp(&s.Cursor, 1)
			case 'J':
				s.moveDown(&s.Cursor, 1)
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
		case tcell.KeyEscape:
			s.Anchor = s.Cursor
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
