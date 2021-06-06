package state

type snapshot struct {
	text           []string
	anchor, cursor cursor
}

type change struct {
	before, after              snapshot
	commonPrefix, commonSuffix int
}

func (s *State) takeSnapshot() {
	s.snapshot.text = make([]string, len(s.Text))
	copy(s.snapshot.text, s.Text)
	s.snapshot.anchor = s.Anchor
	s.snapshot.cursor = s.Cursor
}

func (s *State) recordHistory() {
	var c change

	for c.commonPrefix < len(s.Text) && c.commonPrefix < len(s.snapshot.text) &&
		s.Text[c.commonPrefix] == s.snapshot.text[c.commonPrefix] {
		c.commonPrefix++
	}

	if c.commonPrefix == len(s.Text) && c.commonPrefix == len(s.snapshot.text) {
		return
	}

	for c.commonSuffix < len(s.Text)-c.commonPrefix &&
		c.commonSuffix < len(s.snapshot.text)-c.commonPrefix &&
		s.Text[len(s.Text)-1-c.commonSuffix] ==
			s.snapshot.text[len(s.snapshot.text)-1-c.commonSuffix] {
		c.commonSuffix++
	}

	c.before.text = make(
		[]string,
		len(s.snapshot.text)-c.commonPrefix-c.commonSuffix,
	)
	copy(
		c.before.text,
		s.snapshot.text[c.commonPrefix:len(s.snapshot.text)-c.commonSuffix],
	)
	c.before.anchor = s.snapshot.anchor
	c.before.cursor = s.snapshot.cursor

	c.after.text = make(
		[]string,
		len(s.Text)-c.commonPrefix-c.commonSuffix,
	)
	copy(
		c.after.text,
		s.Text[c.commonPrefix:len(s.Text)-c.commonSuffix],
	)
	c.after.anchor = s.Anchor
	c.after.cursor = s.Cursor

	s.history = s.history[:s.historyHead]
	s.history = append(s.history, c)
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
