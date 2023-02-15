package entity

type DealFieldErr struct {
	Height int64  `bson:"height"`
	Err    string `bson:"err"`
}

func (e DealFieldErr) CollectionName() string {
	return "ex_deal_field_err"
}
