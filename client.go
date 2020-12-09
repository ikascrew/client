package client

import (
	"log"
	"time"

	"github.com/ikascrew/client/config"
	"github.com/ikascrew/client/window"
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

var selector *window.Window

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

		go func() {
			err = virtualController()
			if err != nil {
				log.Println(xerrors.Errorf("virtual Controller error: %w", err))
			}
		}()
		log.Printf("Success Virtual Controller")
	}

	//powermate
	if conf.Powermate {
		pm.HandleFunc(trigger)
		go func() {
			err = pm.Listen("/dev/input/powermate")
			if err != nil {
				log.Printf("powermate listen err: %+v", err)
			} else {
				log.Printf("Success powermate")
			}
		}()
	}

	//Main
	win, err := window.New("ikascrew client", 1536, 768)
	if err != nil {
		return xerrors.Errorf("NewWindow error: %w", err)
	}

	selector = win

	//クライアント描画
	for range time.Tick(10 * time.Millisecond) {
		e := selector.Owner.NextEvent()
		switch e := e.(type) {
		case *window.Part:
			e.Redraw()
		default:
		}
	}

	return nil
}
