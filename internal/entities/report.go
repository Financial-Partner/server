package entities

type Report struct {
	Revenue     int64     `bson:"revenue" json:"revenue"`
	Expenses    int64     `bson:"expenses" json:"expenses"`
	NetProfit   int64     `bson:"net_profit" json:"net_profit"`
	Categories  []string  `bson:"categories" json:"categories"`
	Amounts     []int64   `bson:"amounts" json:"amounts"`
	Percentages []float64 `bson:"percentages" json:"percentages"`
}

type ReportSummary struct {
	Summary string `bson:"summary" json:"summary"`
}
