package handlers

import (
	"github.com/luyuan-li/deal-tx-field/internal/app/config"
	"github.com/luyuan-li/deal-tx-field/internal/app/repository"
	"github.com/luyuan-li/deal-tx-field/libs/pool"
	"io/ioutil"

	"testing"
)

func TestHandlerOneBlock(t *testing.T) {
	data, err := ioutil.ReadFile("/home/lly/GolandProjects/github.com/bianjieai/deal-tx-field/configs/config.toml")
	if err != nil {
		panic(err)
	}
	conf, err := config.ReadConfig(data)
	if err != nil {
		t.Fatal(err.Error())
	}

	block := int64(885460)

	pool.Init(conf)

	repository.InitQMgo(&conf.DataBase)

	InitRouter(conf)

	if err = HandlerOneBlock(block, 0, 0); err != nil {
		t.Fatal(err)
	}
}
