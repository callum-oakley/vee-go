package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/callum-oakley/vee/state"
	"github.com/callum-oakley/vee/ui"
	"github.com/gdamore/tcell/v2"
)

func main() {
	text, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(text), "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	s := state.State{
		FilePath: os.Args[2],
		TabWidth: 4,
		Text:     lines,
	}
	s.Init()

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	r := ui.Renderer{S: &s, Screen: screen}
	r.Render()

	for {
		switch e := screen.PollEvent().(type) {
		case *tcell.EventResize:
			r.Render()
		case *tcell.EventKey:
			if quit := s.HandleKey(e); quit {
				return
			}
			r.Render()
		}
	}
}
