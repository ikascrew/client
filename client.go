package client

import (
	"log"
	"os"

	mc "github.com/ikascrew/core/multicast"
	"github.com/ikascrew/pb"
	pm "github.com/ikascrew/powermate"
	vol "github.com/ikascrew/volumes"
	"github.com/ikascrew/xbox"

	tm "github.com/nsf/termbox-go"
	"golang.org/x/xerrors"
)

var vols *vol.Volumes

func init() {
	vols = vol.New()
	vols.Add("Volume", 300)
	vols.Add("Light", 100)
	vols.Add("Wait", 50)
}

type IkascrewClient struct {
	selector *Window
	testMode bool
}

var ika *IkascrewClient

func Start() error {

	args := os.Args

	ika = &IkascrewClient{}
	if len(args) > 2 {
		ika.testMode = true
	} else {
		ika.testMode = false
	}

	var err error

	cli, err := mc.NewClient()
	if err != nil {
		return xerrors.Errorf("udp client: %w", err)
	}

	log.Println("Search ikascrew Servers")
	acs, err := cli.Find()
	for _, elm := range acs {
		log.Println(elm)
	}

	go func() {
		//XBOX Controller
		xbox.HandleFunc(ika.controller)
		err = xbox.Listen(0)
		if err != nil {
			log.Printf("Controller Listen Error[" + err.Error() + "]")
			return
		}
		log.Printf("Success Controller")
	}()

	go vols.Start()
	//powermate
	go func() {
		pm.HandleFunc(trigger)
		err = pm.Listen("/dev/input/powermate")
		if err != nil {
			log.Printf("powermate Listen Error[" + err.Error() + "]")
			return
		}
		log.Printf("Success powermate")
	}()

	/*
		go func() {
			err = virtualController(ika.controller)
			if err != nil {
				log.Printf("virtual Controller Listen Error[" + err.Error() + "]")
				return
			}
			log.Printf("Success Keyboard")
		}()
	*/

	//Main
	win, err := NewWindow("ikascrew client", 1536, 768)
	if err != nil {
		log.Printf("NewWindow() Error[" + err.Error() + "]")
		return err
	}

	ika.selector = win
	win.SetClient(ika)

	//クライアント描画
	for {
		e := win.window.NextEvent()
		switch e := e.(type) {
		case *Part:
			e.Redraw()
		}
	}

	return err
}

func virtualController(fn func(xbox.Event) error) error {

	//termboxの初期化
	err := tm.Init()
	if err != nil {
		return err
	}
	//プログラム終了時termboxを閉じる
	defer tm.Close()

	xev := xbox.Event{}
	xev.Buttons = make([]bool, 8)
	xev.Axes = make([]int, 2)

	for {
		pev := pm.Event{}
		pev.Type = pm.None

		flag := false
		clearControllerEvent(&xev)

		switch e := tm.PollEvent(); e.Type {
		case tm.EventKey:
			switch e.Key {
			//case tm.KeyEnter:
			case tm.KeyArrowDown:
				pev.Type = pm.Type(pm.Press)
				pev.Value = pm.Value(pm.Down)
			case tm.KeyArrowUp:
				pev.Type = pm.Type(pm.Press)
				pev.Value = pm.Value(pm.Up)
			case tm.KeyArrowRight:
				pev.Type = pm.Type(pm.Rotation)
				pev.Value = pm.Value(pm.Right)
			case tm.KeyArrowLeft:
				pev.Type = pm.Type(pm.Rotation)
				pev.Value = pm.Value(pm.Left)
			case tm.KeyCtrlJ:
				flag = true
				xev.Axes[xbox.CROSS_HORIZONTAL] = -20000
			case tm.KeyCtrlK:
				flag = true
				xev.Axes[xbox.CROSS_HORIZONTAL] = 20000
			case tm.KeyCtrlH:
				flag = true
				xev.Axes[xbox.CROSS_VERTICAL] = -20000
			case tm.KeyCtrlL:
				flag = true
				xev.Axes[xbox.CROSS_VERTICAL] = 20000
			case tm.KeyCtrlQ:
				flag = true
				xev.Buttons[xbox.Y] = true
			case tm.KeyCtrlW:
				flag = true
				xev.Buttons[xbox.X] = true
			case tm.KeyCtrlA:
				flag = true
				xev.Buttons[xbox.B] = true
			case tm.KeyCtrlS:
				flag = true
				xev.Buttons[xbox.A] = true
			default:
			}

			if pev.Type != pm.None {
				err = trigger(pev)
				if err != nil {
					log.Println(err)
				}
			}

			if flag {
				err = fn(xev)
				if err != nil {
					log.Println(err)
				}
			}

		default:
		}
	}

	return nil
}

func trigger(e pm.Event) error {

	val := vols.Get()
	if zero {
		val = 0
		vols.SetCursor(0)
		zero = false
	}

	idx := vols.GetCursor()
	update := false

	switch e.Type {
	case pm.Rotation:
		switch e.Value {
		case pm.Left:
			val -= 2.0
		case pm.Right:
			val += 2.0
		}

		update = true
	case pm.Press:
		switch e.Value {
		case pm.Up:
		case pm.Down:
			idx = idx + 1
			if idx > 2 {
				idx = 0
			}
			vols.SetCursor(idx)
		}
	default:
	}

	if update {

		vols.Set(val)
		var i int64 = int64(idx)

		message := pb.VolumeMessage{
			Index: i,
			Value: val,
		}

		err := ika.callVolume(message)
		if err != nil {
			return err
		}

		//Reply で再度設定
	}

	return nil
}

var zero = false

func setZero() {
	zero = true
}
