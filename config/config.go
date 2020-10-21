package config

import "golang.org/x/xerrors"

var gConf *Config

type ControllerType int

const (
	ControllerTypeNone ControllerType = iota
	ControllerTypeSF
	ControllerTypeJoyCon
)

type Config struct {
	Powermate      bool
	ControllerType ControllerType
}

func (ct ControllerType) String() string {
	switch ct {
	case ControllerTypeNone:
		return "PC Controller"
	case ControllerTypeSF:
		return "SuperFamiryComputer Controller"
	case ControllerTypeJoyCon:
		return "Full Joycon Controller"
	}
	return "Nothing Type"
}

func init() {
	gConf = defaultConfig()
}

func defaultConfig() *Config {

	conf := Config{}

	conf.Powermate = false
	conf.ControllerType = ControllerTypeNone

	return &conf
}

func Get() *Config {
	return gConf
}

func Set(opts ...Option) error {
	for idx, opt := range opts {
		err := opt(gConf)
		if err != nil {
			return xerrors.Errorf("config set error[%d]: %w", idx, err)
		}
	}
	return nil
}
