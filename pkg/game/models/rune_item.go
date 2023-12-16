package models

import "time"

type RuneItem struct {
	Entity     `bson:",inline" json:",inline"`
	MarketInfo `bson:",inline" json:",inline"`

	// Indicates if an item can be dropped in game
	Droppable bool `bson:"droppable" json:"droppable"`

	// How the item changes the wielder's stats
	ItemStats *Modifier `bson:"itemstats" json:"itemstats"`

	// The list of slots to which this rune can be applied to
	ApplicableSlots []int `bson:"applicableslots" json:"applicableslots"`

	// Duration of the potion's effects in seconds.
	//
	// 0 means once active it remains active
	Duration uint `bson:"duration" json:"duration"`

	IsUsed bool `bson:"isused" json:"isused"`

	UsedOn time.Time `bson:"usedon" json:"usedon"`
}

// Implement Type required in the item interface
func (i *RuneItem) Type() int {
	return ItemTypeRune
}

// Implement Id required in the item interface
func (i *RuneItem) Id() string {
	return i.HrId
}

// Implement Modifiers required in the item interface
func (i *RuneItem) Modifiers() *Modifier {
	return i.ItemStats
}

// Implement the attachable interface
func (i *RuneItem) Attach() {

}

// Method to check if the buffs that the rune gives are active. Implement the attachable interface
func (i *RuneItem) IsAttached() bool {
	return false
}

// Method to check if the buffs that the rune gives are active. Implement the attachable interface
func (i *RuneItem) IsDepleted() bool {
	return false
}
