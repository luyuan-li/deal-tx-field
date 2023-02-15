//init client from clientPool.
//client is httpClient of tendermint

package pool

import (
	"context"
	commonPool "github.com/jolestar/go-commons-pool"
	"github.com/luyuan-li/deal-tx-field/internal/app/config"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

var (
	poolObject  *commonPool.ObjectPool
	poolFactory PoolFactory
	ctx         = context.Background()
)

func Init(conf *config.Config) {
	var (
		syncMap sync.Map
	)
	nodeUrls := strings.Split(conf.Server.NodeUrls, ",")
	for _, url := range nodeUrls {
		key := generateId(url)
		endPoint := EndPoint{
			Address:   url,
			Available: true,
		}

		syncMap.Store(key, endPoint)
	}

	poolFactory = PoolFactory{
		peersMap: syncMap,
	}

	config := commonPool.NewDefaultPoolConfig()
	config.MaxTotal = conf.Server.MaxConnectionNum
	config.MaxIdle = conf.Server.InitConnectionNum
	config.MinIdle = conf.Server.InitConnectionNum
	config.TestOnBorrow = true
	config.TestOnCreate = true
	config.TestWhileIdle = true

	poolObject = commonPool.NewObjectPool(ctx, &poolFactory, config)
	poolObject.PreparePool(ctx)
}

// get client from pool
func GetClient() *Client {
	c, err := poolObject.BorrowObject(ctx)
	for err != nil {
		logrus.Errorf("GetClient failed,will try again after 3 seconds :%s", err.Error())
		time.Sleep(3 * time.Second)
		c, err = poolObject.BorrowObject(ctx)
	}

	return c.(*Client)
}

// release client
func (c *Client) Release() {
	err := poolObject.ReturnObject(ctx, c)
	if err != nil {
		logrus.Error(err.Error())
	}
}

func (c *Client) HeartBeat() error {
	http := c.HTTP
	_, err := http.Health(context.Background())
	return err
}

func ClosePool() {
	poolObject.Close(ctx)
}
