package configuration

import (
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Configuration struct {
	config            *viper.Viper
	decoderConfigOpts []viper.DecoderConfigOption
}

func (c *Configuration) Unmarshal(dst any) error {
	err := c.config.ReadInConfig()
	if err != nil {
		return err
	}

	err = c.config.Unmarshal(dst, c.decoderConfigOpts...)
	if err != nil {
		return err
	}

	return nil
}

func New() *Configuration {
	v := viper.NewWithOptions(
		viper.KeyDelimiter("_"),
		viper.EnvKeyReplacer(strings.NewReplacer("-", "_")),
	)

	v.AutomaticEnv()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	decoderConfigOpts := []viper.DecoderConfigOption{
		func(config *mapstructure.DecoderConfig) {
			config.Squash = true
		},
		func(config *mapstructure.DecoderConfig) {
			config.TagName = "config"
		},
	}

	return &Configuration{
		config:            v,
		decoderConfigOpts: decoderConfigOpts,
	}
}
