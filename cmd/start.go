package cmd

import (
	"github.com/luyuan-li/deal-tx-field/internal/app"
	"github.com/luyuan-li/deal-tx-field/internal/app/config"
	"github.com/spf13/cobra"
	"io/ioutil"
)

const (
	defaultLocalConfig = "/home/lly/GolandProjects/github.com/bianjieai/deal-tx-field/configs/config.toml"
)

var (
	localConfig string
	startCmd    = &cobra.Command{
		Use:   "start",
		Short: "Start deal-tx-field Server.",
		Run: func(cmd *cobra.Command, args []string) {
			start()
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&localConfig, "CONFIG", "c", defaultLocalConfig, "conf path: /opt/config.toml")
}

func start() {
	conf := localConf()
	app.Serve(conf)
}

func localConf() *config.Config {
	if localConfig == "" {
		localConfig = defaultLocalConfig
	}
	data, err := ioutil.ReadFile(localConfig)
	if err != nil {
		panic(err)
	}
	conf, err := config.ReadConfig(data)
	if err != nil {
		panic(err)
	}
	return conf
}
