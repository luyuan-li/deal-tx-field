package handlers

import (
	"fmt"
	common_parser "github.com/kaifei-bianjie/common-parser"
	"github.com/kaifei-bianjie/common-parser/codec"
	"github.com/luyuan-li/deal-tx-field/internal/app/config"
	"github.com/luyuan-li/deal-tx-field/internal/app/entity"
	"github.com/luyuan-li/deal-tx-field/internal/app/handler/constant"
	"github.com/luyuan-li/deal-tx-field/internal/app/repository"
	"github.com/luyuan-li/deal-tx-field/internal/app/utils"
	"github.com/luyuan-li/deal-tx-field/libs/msgparser"
	"github.com/luyuan-li/deal-tx-field/libs/pool"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/operator"
	"github.com/sirupsen/logrus"
	types2 "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"strings"
	"sync"
	"time"
)

var (
	_client *pool.Client
	_parser msgparser.MsgParser
)

func InitRouter(conf *config.Config) {
	initBech32Prefix(conf)

	if conf.Server.OnlySupportModule != "" {
		resRouteClient := make(map[string]common_parser.Client, 0)
		modules := strings.Split(conf.Server.OnlySupportModule, ",")
		for _, one := range modules {
			fn, exist := msgparser.RouteClientMap[one]
			if !exist {
				logrus.Fatal("no support module: " + one)
			}
			resRouteClient[one] = fn
		}
		if len(resRouteClient) > 0 {
			msgparser.RouteClientMap = resRouteClient
		}
	}
	_client = pool.GetClient()
	_parser = msgparser.NewMsgParser()
}

func Start() {
	logrus.Infof("handle start... Start_Time:%v", time.Now())

	var (
		minHeight int64
		maxHeight int64
	)

	progress, err := GetLastExecHeight()
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		logrus.Errorf("Start GetLastExecHeight err:%v", err)
		return
	}

	if err == qmgo.ErrNoSuchDocuments {
		minHeight = config.GetConfig().Server.StartHeight
	} else {
		minHeight = progress.StartHeight
	}

	maxHeight = config.GetConfig().Server.EndHeight

	if config.GetConfig().Server.ThreadNum == 0 {
		config.GetConfig().Server.ThreadNum = 1
	}

	token := make(chan struct{}, config.GetConfig().Server.ThreadNum)
	defer close(token)

	var wg sync.WaitGroup

	for i := minHeight; i < maxHeight; i = i + config.GetConfig().Server.IncrHeight {
		wg.Add(1)
		token <- struct{}{}
		startHeight := i
		endHeight := i + config.GetConfig().Server.IncrHeight

		go Handler(&wg, token, startHeight, endHeight)
	}

	wg.Wait()
	logrus.Infof("handle end... End_Time:%v", time.Now())
}

func Handler(wg *sync.WaitGroup, token <-chan struct{}, startHeight, endHeight int64) {
	defer func() {
		<-token
		wg.Done()
	}()

	database := repository.GetClient().Database(config.GetConfig().DataBase.Database)

	var blockList []entity.Block
	query := bson.M{
		"height": bson.M{
			operator.Gt:  startHeight,
			operator.Lte: endHeight,
		},
		"txn": bson.M{
			operator.Gt: 0,
		},
	}
	if err := database.Collection(entity.Block{}.CollectionName()).Find(context.Background(), query).Sort("height").All(&blockList); err != nil {
		logrus.Errorf("Handler Block Find err: %v", err)
		return
	}

	var lastHeight int64

	if len(blockList) > 0 {
		lastHeight = blockList[len(blockList)-1].Height
	}

	queryProgress := bson.M{
		"start_height": startHeight,
		"end_height":   endHeight,
	}

	progress := entity.DealFieldProgress{
		StartHeight:   startHeight,
		EndHeight:     endHeight,
		CurrentHeight: 0,
		LastHeight:    lastHeight,
	}
	_, err := database.Collection(entity.DealFieldProgress{}.CollectionName()).Upsert(context.Background(), queryProgress, progress)
	if err != nil {
		logrus.Errorf("Handler DealFieldProgress Upsert err: %v", err)
		return
	}

	logrus.Infof("start startHeight:%d, endHeight:%d, blockList len:%d", startHeight, endHeight, len(blockList))
	if len(blockList) <= 0 {
		return
	}

	for _, block := range blockList {
		if err := HandlerOneBlock(block.Height, startHeight, endHeight); err != nil {
			log := entity.DealFieldErr{
				Height: block.Height,
				Err:    err.Error(),
			}
			_, err = database.Collection(entity.DealFieldErr{}.CollectionName()).InsertOne(context.Background(), log)
			if err != nil {
				logrus.Errorf("Handler InsertOne err: %v", err)
				continue
			}
		}
	}
	logrus.Infof("end startHeight:%d, endHeight:%d, blockList len:%d", startHeight, endHeight, len(blockList))
}

func HandlerOneBlock(height, startHeight, endHeight int64) error {
	var (
		block *ctypes.ResultBlock
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if v, err := _client.Block(ctx, &height); err != nil {
		time.Sleep(1 * time.Second)
		if v2, err := _client.Block(ctx, &height); err != nil {
			logrus.Errorf("HandlerOneBlock _client Block err :%v", err)
			return err
		} else {
			block = v2
		}
	} else {
		block = v
	}

	if len(block.Block.Data.Txs) <= 0 {
		return nil
	}

	blockResults, err := _client.BlockResults(context.Background(), &height)
	if err != nil {
		time.Sleep(1 * time.Second)
		blockResults, err = _client.BlockResults(context.Background(), &height)
		if err != nil {
			logrus.Errorf("HandlerOneBlock _client BlockResults err :%v", err)
			return err
		}
	}

	if len(block.Block.Txs) != len(blockResults.TxsResults) {
		return err
	}

	if len(block.Block.Txs) > 0 {
		_, err = repository.GetClient().DoTransaction(context.Background(), func(sessCtx context.Context) (interface{}, error) {
			database := repository.GetClient().Database(config.GetConfig().DataBase.Database)
			for i, txBytes := range block.Block.Txs {
				var (
					updateFeeGranter string
					updateFeePayer   string
					updateFeeGrantee string
				)

				txResult := blockResults.TxsResults[i]

				txHash := utils.BuildHex(txBytes.Hash())

				authTx, err := codec.GetSigningTx(txBytes)
				if err != nil {
					logrus.Errorf("HandlerOneBlock SyncTx UpdateOne is err, height:%d,tx_hash:%s,err:%v", block.Block.Height, txHash, err)
					return nil, err
				}

				feeGranter := authTx.FeeGranter()
				if feeGranter != nil {
					updateFeeGranter = feeGranter.String()
					updateFeePayer = feeGranter.String()
				} else {

					feePayer := authTx.FeePayer()
					if feePayer != nil {
						updateFeePayer = feePayer.String()
					} else {
						if authTx.GetSigners() != nil {
							updateFeePayer = authTx.GetSigners()[0].String()
						}
					}
				}

				feeGrantee := GetFeeGranteeFromEvents(txResult.Events, constant.Grantee)
				if feeGrantee != "" {
					updateFeeGrantee = feeGrantee
				}

				queryTx := bson.M{
					"height":  block.Block.Height,
					"tx_hash": txHash,
				}

				updateTx := bson.M{
					operator.Set: bson.M{
						"fee_granter": updateFeeGranter,
						"fee_payer":   updateFeePayer,
						"fee_grantee": updateFeeGrantee,
					},
				}

				if err = database.Collection(entity.SyncTx{}.CollectionName()).UpdateOne(sessCtx, queryTx, updateTx); err != nil {
					logrus.Errorf("HandlerOneBlock SyncTx UpdateOne is err, height:%d,tx_hash:%s,err:%v", block.Block.Height, txHash, err)
					return nil, err
				}

				queryProgress := bson.M{
					"start_height": startHeight,
					"end_height":   endHeight,
				}

				updateProgress := bson.M{
					operator.Set: bson.M{
						"current_height": block.Block.Height,
					},
				}

				if err = database.Collection(entity.DealFieldProgress{}.CollectionName()).UpdateOne(sessCtx, queryProgress, updateProgress); err != nil {
					logrus.Errorf("HandlerOneBlock DealFieldProgress UpdateOne is err, height:%d,tx_hash:%s,err:%v", block.Block.Height, txHash, err)
					return nil, err
				}

			}

			return nil, nil
		})

	}

	return err
}

func GetFeeGranteeFromEvents(events []types2.Event, key string) string {
	for _, val := range events {
		if val.Type == constant.UseGrantee /*|| val.Type == constant.SetGrantee*/ {
			for _, attribute := range val.Attributes {
				if fmt.Sprintf("%s", attribute.Key) == key {
					return fmt.Sprintf("%s", attribute.Value)
				}
			}
		}
	}
	return ""
}

func GetLastExecHeight() (entity.DealFieldProgress, error) {
	database := repository.GetClient().Database(config.GetConfig().DataBase.Database)
	var progress entity.DealFieldProgress

	bsonProject := bson.D{{operator.Project, bson.D{
		{"start_height", 1},
		{"end_height", 1},
		{"current_height", 1},
		{"last_height", 1},
		{"diff", bson.D{{operator.Ne, bson.A{"$current_height", "$last_height"}}}},
	}}}
	bsonMatch := bson.D{{operator.Match, bson.D{
		{"diff", true},
	}}}

	bsonSort := bson.D{{operator.Sort, bson.D{
		{"current_height", 1},
	}}}

	bsonLimit := bson.D{{operator.Limit, 1}}

	err := database.Collection(entity.DealFieldProgress{}.CollectionName()).Aggregate(context.Background(), mongo.Pipeline{bsonProject, bsonMatch, bsonSort, bsonLimit}).One(&progress)
	if err != nil && err != qmgo.ErrNoSuchDocuments {
		logrus.Errorf("Start DealFieldProgress Find err:%v", err)
		return progress, err
	}
	return progress, err
}
