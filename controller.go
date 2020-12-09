package client

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ikascrew/xbox"
	"golang.org/x/xerrors"
)

var (
	JoyconButtons = []string{"A", "B", "X", "Y", "L", "R", "BACK", "START", "L_JOY", "R_JOY"}
	JoyconAxes    = []string{"LEFT_JOY_H", "LEFT_JOY_V", "ZLR", "RIGHT_JOY_H", "RIGHT_JOY_V", "CROSS_H", "CROSS_V"}
)

type Controller struct {
	body *xbox.Controller
}

func createController(id int) (*Controller, error) {

	ctrl, err := xbox.New(id,
		xbox.Logger(log.New(os.Stdout, "[CTRL]", log.LstdFlags|log.Lshortfile)),
		xbox.Duration(40),
		xbox.AxisMargin(3000),
	)
	if err != nil {
		return nil, xerrors.Errorf("controller error: %w", err)
	}

	if true {
		err = ctrl.ButtonNames(JoyconButtons...)
		if err != nil {
			return nil, xerrors.Errorf("Button error(10): %w", err)
		}

		err = ctrl.AxisNames(JoyconAxes...)
		if err != nil {
			return nil, xerrors.Errorf("JoyStick(7): %w", err)
		}
	} else {
		err = ctrl.ButtonNames("A", "B", "X", "Y", "L", "R", "BACK", "START")
		if err != nil {
			return nil, xerrors.Errorf("Button error(8): %w", err)
		}
		err = ctrl.AxisNames("CROSS_H", "CROSS_V")
		if err != nil {
			return nil, xerrors.Errorf("JoyStick(2): %w", err)
		}
	}

	rtn := Controller{
		body: ctrl,
	}

	return &rtn, nil
}

func (ctrl *Controller) Listen() error {
	ch := ctrl.body.Event()
	for {
		select {
		case ev := <-ch:
			if ev.Error() != nil {
				return xerrors.Errorf("controller error: %w", ev.Error())
			} else {
				raise(ev)
			}
		default:
		}

		if ctrl.body.Closed() {
			break
		}
	}
	return nil
}

type Event struct {
	Type  EventType
	Value int
}

type EventType int

const (
	EventList EventType = iota
	EventNext
	EventUpper
	EventLowwer
	EventSelectList
	EventDeleteNext
	EventSelectNext
	EventView
	EventSync
	EventNone
)

func createEvent(e *xbox.Event) *Event {

	ev := Event{}
	t := EventNone
	v := 0

	for _, btn := range e.Buttons {
		switch btn.Name {
		case "A":
			t = EventSelectList
		case "B":
			t = EventDeleteNext
		case "X":
			t = EventSelectNext
		case "Y":
			t = EventView
		case "L":
			t = EventUpper
		case "R":
			t = EventLowwer
		case "START":
			t = EventSync
		}
	}

	for _, axis := range e.Axes {
		//"LEFT_JOY_H", "LEFT_JOY_V", "ZLR", "RIGHT_JOY_V", "RIGHT_JOY_H", "CROSS_H", "CROSS_V"
		switch axis.Name {
		case "RIGHT_JOY_H":
			t = EventNext
			v = axis.Value
		case "LEFT_JOY_V":
			t = EventList
			v = axis.Value
		}
	}

	//TODO Superfamicon

	ev.Type = t
	ev.Value = v

	return &ev
}

func raise(e *xbox.Event) error {

	ev := createEvent(e)

	switch ev.Type {
	case EventSync:
		callSync()
	case EventList:
		selector.List.SetCursor(ev.Value / 2)
		selector.List.Push()
	case EventNext:
		selector.Next.SetCursor(ev.Value / 2)
		selector.Next.Push()
	case EventUpper:
		selector.List.ZeroCursor()
	case EventLowwer:
		selector.List.MaxCursor()
	case EventSelectList:
		res := selector.List.Get()
		if res != "" {
			err := selector.Next.Add(res)
			if err != nil {
				// TODO 無視
			}
			selector.Next.Push()
		} else {
			log.Printf("Selector Error:" + "No Index")
		}
	case EventDeleteNext:
		err := selector.Next.Delete()
		if err != nil {
			// TODO 無視
			log.Printf("Pusher Delete Error:", err)
		}
		selector.Next.Push()
	case EventSelectNext:
		res := selector.Next.Get()
		if res != "" {

			t := "file"
			if strings.Index(res, "countdown") >= 0 {
				t = "countdown"
			}

			idx := strings.LastIndex(res, "/")
			buf := ""
			if idx != -1 {
				buf = res[idx+1:]
				buf = strings.Replace(buf, ".jpg", "", -1)
			}

			id, err := strconv.Atoi(buf)
			if err != nil {
				return fmt.Errorf("Efffect id error:[%s]", res)
			}

			err = callEffect(int64(id), t)
			if err != nil {
				log.Printf("callEffect[%+v]", err)
			} else {

				//0
				setZero()

				selector.Next.Delete()
				selector.Next.Push()
			}
		} else {
			log.Printf("Pusher Error: No Index")
		}
	case EventView:
		res := selector.List.Get()
		if res != "" {

			err := selector.Player.SetFile(res)
			if err != nil {
				log.Printf("%+v", err)
			} else {

				selector.Player.Draw()
				selector.Player.Push()
			}

		} else {
			log.Printf("Pusher Error: No Index")
		}
	}

	return nil
}
