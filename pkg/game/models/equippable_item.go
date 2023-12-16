package models

type EquippableItem struct {
	Uid        string `bson:"uid" json:"Uid"`
	Entity     `bson:",inline" json:",inline"`
	MarketInfo `bson:",inline" json:",inline"`

	// Indicates if an item can be dropped in game
	Droppable bool `bson:"droppable" json:"-"`

	// can an item be buffed by a rune
	Buffable bool `bson:"buffable" json:"buffable"`

	// Which slot the item occupies when equipped
	Slot int `bson:"slot" json:"slot"`

	// Rarity of an item
	Rarity int `bson:"rarity" json:"rarity"`

	// Minimum level required to equip item
	MinLevel int `bson:"minlevel" json:"minlevel"`

	// On which level should we stop dropping this item.
	StopDropLevel int `bson:"stopdrop" json:"-"`

	// How the item changes the wielder's stats
	ItemStats *Modifier `bson:"itemstats" json:"itemstats"`

	Runes []*RuneItem `bson:"runes" json:"runes"`
}

// Implement Type required in the item interface
func (i *EquippableItem) Type() int {
	return ItemTypeEquippable
}

// Implement Id required in the item interface
func (i *EquippableItem) Id() string {
	return i.HrId + "-" + i.Uid
}

// Implement Modifiers required in the item interface
func (i *EquippableItem) Modifiers() *Modifier {
	return i.ItemStats
}

// Implement Equip required in the Equippable interface
func (i *EquippableItem) Equip(_c *Character) {
	_c.EquippedItems[i.Slot] = i.Id()
}

// Implement Unequip required in the Equippable interface
func (i *EquippableItem) Unequip(_c *Character) {
	_c.EquippedItems[i.Slot] = ""
}

// Implement AddRune required in the Equippable interface
func (i *EquippableItem) AddRune(_rune *RuneItem) {
	i.Runes = append(i.Runes, _rune)
}

// Implement AddRune required in the Equippable interface
func (i *EquippableItem) AttachedRunes() []*RuneItem {
	return i.Runes
}
