package state

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	rw "github.com/mattn/go-runewidth"
)

type State struct {
	FilePath string
	TabWidth int
	Text     []string
	Cursor   struct{ X, Y, col int }
	Msg      string
}

func visualWidth(col, tabWidth int, char rune) int {
	if char == '\t' {
		return tabWidth - col%tabWidth
	}
	return rw.RuneWidth(char)
}

func (s *State) updateCol() {
	s.Cursor.col = 0
	for _, char := range s.Text[s.Cursor.Y][:s.Cursor.X] {
		s.Cursor.col += visualWidth(s.Cursor.col, s.TabWidth, char)
	}
}

func (s *State) updateX() {
	s.Cursor.X = -1
	col := 0
	var char rune
	for s.Cursor.X, char = range s.Text[s.Cursor.Y] {
		col += visualWidth(col, s.TabWidth, char)
		if col > s.Cursor.col {
			break
		}
	}
}

func (s *State) MoveLeft() {
	if s.Cursor.X <= 0 {
		return
	}
	_, size := utf8.DecodeLastRuneInString(s.Text[s.Cursor.Y][:s.Cursor.X])
	s.Cursor.X -= size
	s.updateCol()
}

func (s *State) MoveRight() {
	if len(s.Text[s.Cursor.Y]) == 0 {
		return
	}
	_, sizeLast := utf8.DecodeLastRuneInString(s.Text[s.Cursor.Y])
	if s.Cursor.X == len(s.Text[s.Cursor.Y])-sizeLast {
		return
	}
	_, size := utf8.DecodeRuneInString(s.Text[s.Cursor.Y][s.Cursor.X:])
	s.Cursor.X += size
	s.updateCol()
}

func (s *State) MoveUp(n int) {
	if s.Cursor.Y == 0 {
		return
	}
	s.Cursor.Y -= n
	if s.Cursor.Y < 0 {
		s.Cursor.Y = 0
	}
	s.updateX()
}

func (s *State) MoveDown(n int) {
	if s.Cursor.Y == len(s.Text)-1 {
		return
	}
	s.Cursor.Y += n
	if s.Cursor.Y > len(s.Text)-1 {
		s.Cursor.Y = len(s.Text) - 1
	}
	s.updateX()
}

func (s *State) MoveStartOfLine() {
	if len(s.Text[s.Cursor.Y]) == 0 {
		return
	}
	var char rune
	for s.Cursor.X, char = range s.Text[s.Cursor.Y] {
		if char != ' ' && char != '\t' {
			break
		}
	}
	s.updateCol()
}

func (s *State) MoveEndOfLine() {
	if len(s.Text[s.Cursor.Y]) == 0 {
		return
	}
	_, sizeLast := utf8.DecodeLastRuneInString(s.Text[s.Cursor.Y])
	s.Cursor.X = len(s.Text[s.Cursor.Y]) - sizeLast
	s.updateCol()
}

var reStartOfWord = regexp.MustCompile(
	`([[:word:]]+|[[:punct:]]+)[[:blank:]]*$`,
)

func (s *State) MoveStartOfWord() {
	if len(s.Text[s.Cursor.Y]) == 0 {
		return
	}
	if match := reStartOfWord.FindStringIndex(
		s.Text[s.Cursor.Y][:s.Cursor.X],
	); match != nil {
		s.Cursor.X = match[0]
	}
}

var reEndOfWord = regexp.MustCompile(
	`^[[:blank:]]*([[:word:]]+|[[:punct:]]+)`,
)

func (s *State) MoveEndOfWord() {
	if len(s.Text[s.Cursor.Y]) == 0 {
		return
	}
	if match := reEndOfWord.FindStringIndex(
		s.Text[s.Cursor.Y][s.Cursor.X+1:],
	); match != nil {
		s.Cursor.X += match[1]
	}
}

func (s *State) DebugUnicode() {
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
