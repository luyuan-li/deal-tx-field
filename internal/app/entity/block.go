package entity

type Block struct {
	Height   int64  `bson:"height"`
	Hash     string `bson:"hash"`
	Txn      int64  `bson:"txn"`
	Time     int64  `bson:"time"`
	Proposer string `bson:"proposer"`
}

func (d Block) CollectionName() string {
	return "sync_block"
}
