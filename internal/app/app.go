package app

import (
	"github.com/luyuan-li/deal-tx-field/internal/app/config"
	handlers "github.com/luyuan-li/deal-tx-field/internal/app/handler"
	"github.com/luyuan-li/deal-tx-field/internal/app/repository"
	"github.com/luyuan-li/deal-tx-field/libs/pool"
)

func Serve(conf *config.Config) {
	pool.Init(conf)

	repository.InitQMgo(&conf.DataBase)

	handlers.InitRouter(conf)

	handlers.Start()
}
