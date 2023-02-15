package repository

import (
	"context"
	"fmt"
	"github.com/luyuan-li/deal-tx-field/internal/app/config"

	"github.com/qiniu/qmgo"
)

var mgoCli *qmgo.Client

func InitQMgo(config *config.DataBaseConf) {
	var maxPoolSize uint64 = 4096
	cli, err := qmgo.NewClient(context.Background(), &qmgo.Config{
		Uri:         config.Url,
		MaxPoolSize: &maxPoolSize,
	})
	if err != nil {
		fmt.Println(err)
	}
	mgoCli = cli
}

func GetClient() *qmgo.Client {
	return mgoCli
}

func MongoDbStatus() bool {
	err := mgoCli.Ping(5)
	if err != nil {
		return false
	}
	return true
}
