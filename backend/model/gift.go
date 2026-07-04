package model

type GiftConfig struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

var GiftList = []GiftConfig{
	{Name: "heart", Score: 10},
	{Name: "rocket", Score: 100},
}

func GetGiftScore(giftType string) int {
	for _, g := range GiftList {
		if g.Name == giftType {
			return g.Score
		}
	}
	return 0
}

func IsValidGift(giftType string) bool {
	return GetGiftScore(giftType) > 0
}
