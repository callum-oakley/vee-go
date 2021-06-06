package state

import (
	"regexp"
	"unicode/utf8"

	rw "github.com/mattn/go-runewidth"
)

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
