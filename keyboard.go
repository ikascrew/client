package client

import (
	"fmt"

	pm "github.com/ikascrew/powermate"
	"github.com/ikascrew/xbox"
	tm "github.com/nsf/termbox-go"
	"golang.org/x/xerrors"
)

func virtualController() error {

	err := tm.Init()
	if err != nil {
		return xerrors.Errorf("termbox init: %w", err)
	}
	defer tm.Close()

	for {

		var tmEv *xbox.Event
		var pmEv *pm.Event

		switch e := tm.PollEvent(); e.Type {
		case tm.EventKey:
			switch e.Key {
			case tm.KeyEnter:
				tmEv = newVirtualEvent()
				btn := xbox.NewButton(0, "A")
				tmEv.Buttons = append(tmEv.Buttons, btn)
			case tm.KeySpace:
				tmEv = newVirtualEvent()
				btn := xbox.NewButton(0, "X")
				tmEv.Buttons = append(tmEv.Buttons, btn)
			case tm.KeyArrowDown:
			case tm.KeyArrowUp:
			case tm.KeyArrowRight:
				pmEv = new(pm.Event)
				pmEv.Type = pm.Rotation
				pmEv.Value = pm.Right
			case tm.KeyArrowLeft:
				pmEv = new(pm.Event)
				pmEv.Type = pm.Rotation
				pmEv.Value = pm.Left
			case tm.KeyCtrlJ:
				tmEv = newVirtualEvent()
				tmEv.Axes = append(tmEv.Axes, xbox.NewAxis(1, "LEFT_JOY_V", 20000))
			case tm.KeyCtrlK:
				tmEv = newVirtualEvent()
				tmEv.Axes = append(tmEv.Axes, xbox.NewAxis(1, "LEFT_JOY_V", -20000))
			case tm.KeyCtrlH:
				tmEv = newVirtualEvent()
				tmEv.Axes = append(tmEv.Axes, xbox.NewAxis(3, "RIGHT_JOY_H", -70000))
			case tm.KeyCtrlL:
				tmEv = newVirtualEvent()
				tmEv.Axes = append(tmEv.Axes, xbox.NewAxis(3, "RIGHT_JOY_H", 70000))
			case tm.KeyCtrlQ:
			case tm.KeyCtrlC:
				fmt.Println("Window Bye!")
				return nil
			case tm.KeyCtrlW:
			case tm.KeyCtrlA:
			case tm.KeyCtrlS:
			default:
			}

			if tmEv != nil {
				raise(tmEv)
			} else if pmEv != nil {
				trigger(*pmEv)
			}

		default:
		}
	}

	return nil
}

func newVirtualEvent() *xbox.Event {
	ev := xbox.Event{}
	ev.Buttons = make([]*xbox.Button, 0, len(JoyconButtons))
	ev.Axes = make([]*xbox.Axis, 0, len(JoyconAxes))
	return &ev
}
