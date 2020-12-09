package client

import (
	"github.com/ikascrew/pb"
	pm "github.com/ikascrew/powermate"
)

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
