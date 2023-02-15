package config

import (
	"bytes"
	"github.com/spf13/viper"
)

var (
	_config *Config
)

type (
	Config struct {
		DataBase DataBaseConf `mapstructure:"database"`
		Server   ServerConf   `mapstructure:"server"`
	}

	DataBaseConf struct {
		Url      string `mapstructure:"url"`
		Database string `mapstructure:"database"`
	}

	ServerConf struct {
		NodeUrls          string `mapstructure:"node_urls"`
		ThreadNum         int    `mapstructure:"thread_num"`
		StartHeight       int64  `mapstructure:"start_height"`
		EndHeight         int64  `mapstructure:"end_height"`
		IncrHeight        int64  `mapstructure:"incr_height"`
		MaxConnectionNum  int    `mapstructure:"max_connection_num"`
		InitConnectionNum int    `mapstructure:"init_connection_num"`
		Bech32AccPrefix   string `mapstructure:"bech32_acc_prefix"`
		OnlySupportModule string `mapstructure:"only_support_module"`
	}
)

func ReadConfig(data []byte) (*Config, error) {
	v := viper.New()
	v.SetConfigType("toml")
	reader := bytes.NewReader(data)
	err := v.ReadConfig(reader)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}

	_config = &conf
	return &conf, nil
}

func GetConfig() *Config {
	return _config
}
