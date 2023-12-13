package models

type MarketInfo struct {
	// Indicates if the item can be sold or auctioned
	IsMarketable bool `bson:"marketable" json:"marketable"`

	// Indicates if the item can be traded to another
	IsTradable bool `bson:"tradable" json:"tradable"`

	// Gold price required to buy the entity
	BuyCostGold int `bson:"buycostgold" json:"buycostgold"`

	// Gems price required to buy the entity
	BuyCostGems int `bson:"buycostgems" json:"buycostgems"`

	// Gold price required to sell the entity
	SellCostGold int `bson:"sellcostgold" json:"sellcostgold"`

	// Gems price required to sell the entity
	SellCostGems int `bson:"sellcostgems" json:"sellcostgems"`
}
