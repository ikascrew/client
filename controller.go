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
	JoyconAxes    = []string{"LEFT_JOY_H", "LEFT_JOY_V", "ZLR", "RIGHT_JOY_V", "RIGHT_JOY_H", "CROSS_H", "CROSS_V"}
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
		}
	}

	for _, axis := range e.Axes {
		//"LEFT_JOY_H", "LEFT_JOY_V", "ZLR", "RIGHT_JOY_V", "RIGHT_JOY_H", "CROSS_H", "CROSS_V"
		switch axis.Name {
		case "LEFT_JOY_H":
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
	case EventList:
		selector.list.setCursor(ev.Value / 2)
		selector.list.Push()
	case EventNext:
		selector.next.setCursor(ev.Value)
		selector.next.Push()
	case EventUpper:
		selector.list.zeroCursor()
	case EventLowwer:
		selector.list.maxCursor()
	case EventSelectList:
		res := selector.list.get()
		if res != "" {
			err := selector.next.add(res)
			if err != nil {
				// TODO 無視
			}
			selector.next.Push()
		} else {
			log.Printf("Selector Error:" + "No Index")
		}
	case EventDeleteNext:
		err := selector.next.delete()
		if err != nil {
			// TODO 無視
			log.Printf("Pusher Delete Error:", err)
		}
		selector.next.Push()
	case EventSelectNext:
		res := selector.next.get()
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
				fmt.Println(buf)
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

				selector.next.delete()
				selector.next.Push()
			}
		} else {
			log.Printf("Pusher Error: No Index")
		}
	case EventView:
		res := selector.list.get()
		if res != "" {

			selector.player.setFile(res)
			selector.player.Draw()
			selector.player.Push()

		} else {
			log.Printf("Pusher Error: No Index")
		}
	}

	return nil
}
