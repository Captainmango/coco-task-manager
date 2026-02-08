package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var Config config

type config struct {
	CrontabFile  string `env:"CRONTAB_FILE" envDefault:"./e2e/storage/crontab"`
	RabbitMQHost string `env:"RABBITMQ_HOST" envDefault:"localhost:5672/"`
}

type ConfigOptFn func(o *opts)
type opts struct {
	useDotEnv  bool
	dotEnvPath string
}

func WithDotEnv() ConfigOptFn {
	return func(o *opts) {
		o.useDotEnv = true
	}
}

func WithDotEnvPath(path string) ConfigOptFn {
	return func(o *opts) {
		if !o.useDotEnv {
			o.useDotEnv = true
		}

		o.dotEnvPath = path
	}
}

func BootstrapConfig(configOpts ...ConfigOptFn) {
	opts := &opts{}

	for _, optFn := range configOpts {
		optFn(opts)
	}

	if opts.useDotEnv {
		if opts.dotEnvPath != "" {
			//nolint:errcheck // We don't care if this blows up.
			godotenv.Load(opts.dotEnvPath)
		} else {
			//nolint:errcheck // ... Or this
			godotenv.Load()
		}
	}

	err := env.Parse(&Config)

	if err != nil {
		fmt.Println(err)
		return
	}
}
