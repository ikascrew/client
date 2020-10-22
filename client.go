package client

import (
	"log"

	"github.com/ikascrew/client/config"
	"github.com/ikascrew/pb"
	pm "github.com/ikascrew/powermate"
	vol "github.com/ikascrew/volumes"

	"golang.org/x/xerrors"
)

var vols *vol.Volumes

func init() {
	vols = vol.New()
	vols.Add("Volume", 300)
	vols.Add("Light", 100)
	vols.Add("Wait", 50)
}

var selector *Window

func Start(opts ...config.Option) error {

	err := config.Set(opts...)
	if err != nil {
		return xerrors.Errorf("option set error: %w", err)
	}

	/*

		cli, err := mc.NewClient()
		if err != nil {
			return xerrors.Errorf("udp client: %w", err)
		}

			    // TODO: windows multicast support??
				acs, err := cli.Find()
				if err != nil {
					return xerrors.Errorf("not found server: %w", err)
				}

				for _, elm := range acs {
					log.Println(elm)
				}
	*/

	log.Println("ServerVolume Start")
	go vols.Start()

	conf := config.Get()

	if conf.ControllerType != config.ControllerTypeNone {

		ctrl, err := createController(0)
		if err != nil {
			return xerrors.Errorf("createController Error: %w", err)
		}

		go func() {
			err = ctrl.Listen()
			log.Printf("Listen error: %+v\n", err)
		}()

	} else {
		//TODO キーボードのみで設定を可能にする
		//err = virtualController(ika.controller)
		//if err != nil {
		//}
		//log.Printf("Success Virtual Controller")
		return xerrors.Errorf("virtual Controller not supported.")
	}

	//powermate
	if conf.Powermate {
		pm.HandleFunc(trigger)
		err = pm.Listen("/dev/input/powermate")
		if err != nil {
			return xerrors.Errorf("powermate Listen Error : %w", err)
		}
		log.Printf("Success powermate")
	}

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
		return xerrors.Errorf("NewWindow error: %w", err)
	}

	selector = win

	//クライアント描画
	for {
		e := selector.window.NextEvent()
		switch e := e.(type) {
		case *Part:
			e.Redraw()
		default:
		}
	}

	return nil
}

/*
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
*/

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

		err := callVolume(message)
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
