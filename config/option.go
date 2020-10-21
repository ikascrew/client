package config

type Option func(*Config) error

func UsePowermate() Option {
	return func(conf *Config) error {
		conf.Powermate = true
		return nil
	}
}

func Controller(ct ControllerType) Option {
	return func(conf *Config) error {
		conf.ControllerType = ct
		return nil
	}
}
