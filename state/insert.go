package state

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
