package state

import (
	"fmt"
	"os"
)

func (s *State) save() {
	f, err := os.Create(s.FilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, line := range s.Text {
		fmt.Fprintln(f, line)
	}
}
