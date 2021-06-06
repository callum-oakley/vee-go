package state

import (
	"fmt"
	"unicode/utf8"

	rw "github.com/mattn/go-runewidth"
)

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
