package domain

type Recommendation struct {
	Ticker     string  `gorm:"primaryKey" json:"ticker"`
	Brokerage  string  `gorm:"primaryKey" json:"brokerage"`
	Company    string  `json:"company"`
	Action     string  `json:"action"`
	RatingFrom string  `json:"rating_from"`
	RatingTo   string  `json:"rating_to"`
	TargetFrom float64 `json:"target_from"`
	TargetTo   float64 `json:"target_to"`
}
