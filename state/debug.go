package state

import (
	"fmt"
	"os"
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

func (s *State) debugLog(msg string) {
	f, err := os.OpenFile("vee.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%v\n", msg)); err != nil {
		panic(err)
	}
}
