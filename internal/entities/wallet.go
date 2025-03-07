package entities

type Wallet struct {
	Diamonds int64 `bson:"diamonds" json:"diamonds"`
	Savings  int64 `bson:"savings" json:"savings"`
}
