package models

type ConsumableItem struct {
	Entity     `bson:",inline" json:",inline"`
	MarketInfo `bson:",inline" json:",inline"`

	// Indicates if an item can be dropped in game
	Droppable bool `bson:"droppable" json:"droppable"`

	// How the item changes the wielder's stats
	ItemStats *Modifier `bson:"itemstats" json:"itemstats"`

	// Duration of the potion's effects in seconds.
	//
	// 0 means once active it remains active
	Duration uint `bson:"duration" json:"duration"`
}

// Implement Type required in the item interface
func (i *ConsumableItem) Type() int {
	return ItemTypeConsumable
}

// Implement Id required in the item interface
func (i *ConsumableItem) Id() string {
	return i.HrId
}

// Implement Modifiers required in the item interface
func (i *ConsumableItem) Modifiers() *Modifier {
	return i.ItemStats
}

// Implement Modifiers required in the consumable interface
func (i *ConsumableItem) Consume() {
	// TODO: implement
}
