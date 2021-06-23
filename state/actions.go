package state

func (s *State) newLineAbove() {
	s.applyDiff(diff{
		start:  s.Cursor.Y,
		before: []string{},
		after:  []string{""},
	})
	s.setCursorY(&s.Cursor, s.Cursor.Y)
	s.Anchor = s.Cursor
}

func (s *State) normalisedSelection() (cursor, cursor) {
	if s.Cursor.Y < s.Anchor.Y ||
		s.Cursor.Y == s.Anchor.Y && s.Cursor.X < s.Anchor.X {
		return s.Cursor, s.Anchor
	}
	return s.Anchor, s.Cursor
}

func (s *State) normaliseSelection() {
	s.Anchor, s.Cursor = s.normalisedSelection()
}

func (s *State) delete() {
	s.normaliseSelection()
	var after []string
	if s.Anchor.X == -1 && s.Cursor.X == -1 {
		after = []string{}
	} else if s.Anchor.X == -1 {
		after = []string{s.Text[s.Cursor.Y][s.xRightOf(&s.Cursor):]}
	} else if s.Cursor.X == -1 {
		after = []string{s.Text[s.Anchor.Y][:s.Anchor.X]}
	} else {
		after = []string{s.Text[s.Anchor.Y][:s.Anchor.X] +
			s.Text[s.Cursor.Y][s.xRightOf(&s.Cursor):]}
	}
	s.applyDiff(diff{
		start:  s.Anchor.Y,
		before: s.Text[s.Anchor.Y : s.Cursor.Y+1],
		after:  after,
	})
	s.setCursorY(&s.Anchor, s.Anchor.Y)
	s.Cursor = s.Anchor
}

func (s *State) deleteLines() {
	s.normaliseSelection()
	s.applyDiff(diff{
		start:  s.Anchor.Y,
		before: s.Text[s.Anchor.Y : s.Cursor.Y+1],
		after:  []string{},
	})
	if s.Anchor.Y >= len(s.Text) {
		s.Anchor.Y = len(s.Text) - 1
	}
	s.setCursorY(&s.Anchor, s.Anchor.Y)
	s.Cursor = s.Anchor
}
