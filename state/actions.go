package state

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
