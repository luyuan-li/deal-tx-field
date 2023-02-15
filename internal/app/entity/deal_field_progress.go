package entity

type DealFieldProgress struct {
	StartHeight   int64 `bson:"start_height"`
	EndHeight     int64 `bson:"end_height"`
	CurrentHeight int64 `bson:"current_height"`
	LastHeight    int64 `bson:"last_height"`
}

func (e DealFieldProgress) CollectionName() string {
	return "ex_deal_field_progress"
}
