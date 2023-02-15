package entity

type (
	SyncTx struct {
		Time          int64      `bson:"time"`
		Height        int64      `bson:"height"`
		TxHash        string     `bson:"tx_hash"`
		Type          string     `bson:"type"`
		Memo          string     `bson:"memo"`
		Status        int64      `bson:"status"`
		Log           string     `bson:"log"`
		Fee           FeeAmount  `bson:"fee"`
		Signers       []string   `bson:"signers"`
		Msgs          []DocTxMsg `bson:"msgs"`
		TxId          int64      `bson:"tx_id"`
		Addrs         []string   `bson:"addrs"`
		ContractAddrs []string   `bson:"contract_addrs"`
		EventsNew     []EventNew `bson:"events_new"`
	}

	Event struct {
		Type       string   `bson:"type"`
		Attributes []KvPair `bson:"attributes"`
	}

	KvPair struct {
		Key   string `bson:"key"`
		Value string `bson:"value"`
	}

	EventNew struct {
		MsgIndex uint32  `bson:"msg_index" json:"msg_index"`
		Events   []Event `bson:"events"`
	}
)

func (e SyncTx) CollectionName() string {
	return "sync_tx"
}

type SyncTxs []SyncTx

type MtMsg struct {
	Id        string `bson:"id" json:"id"`
	Name      string `bson:"name" json:"name"`
	DenomId   string `bson:"denom_id" json:"denom_id"`
	Amount    string `bson:"amount" json:"amount"`
	Data      string `bson:"data" json:"data"`
	Sender    string `bson:"sender" json:"sender"`
	Recipient string `bson:"recipient" json:"recipient"`
}

type MtDenomMsg struct {
	Id        string `bson:"id" json:"id"`
	Name      string `bson:"name" json:"name"`
	Data      string `bson:"data" json:"data"`
	Sender    string `bson:"sender" json:"sender"`
	Recipient string `bson:"recipient" json:"recipient"`
}

type FeeAmount struct {
	Amount []DenomAmount `bson:"amount" json:"amount"`
	Gas    int64         `bson:"gas" json:"gas"`
}

type DocTxMsg struct {
	Type string                 `bson:"type"`
	Msg  map[string]interface{} `bson:"msg"`
}

type DenomAmount struct {
	Denom  string `bson:"denom" json:"denom"`
	Amount string `bson:"amount" json:"amount"`
}
