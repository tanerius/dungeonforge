package models

// Item interface for all items
type Item interface {
	Id() string
	Modifiers() *Modifier
}

// Consumable interface for consumable items
type Consumable interface {
	Item
	Consume()
}

// Equippable interface for equippable items
type Equippable interface {
	Item
	Equip()
	AddRune(rune *RuneItem)
}

// RuneAttachable interface for items that runes can be attached to
type RuneAttachable interface {
	AddRune(rune *RuneItem)
}
