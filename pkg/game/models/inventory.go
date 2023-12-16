package models

const (
	ItemTypeUnknown    int = 0
	ItemTypeEquippable int = 1
	ItemTypeConsumable int = 2
	ItemTypeRune       int = 3
)

// Item interface for all items
type Item interface {
	Id() string
	Type() int
	Modifiers() *Modifier
}

// Consumable interface for consumable items
type Consumable interface {
	Item
	Consume(*Character)
}

// Equippable interface for equippable items
type Equippable interface {
	Item
	Equip(*Character)
	Unequip(*Character)
	AddRune(*RuneItem)

	AttachedRunes() []*RuneItem
}

// Attachable interface for items that can attach to other items
type Attachable interface {
	Item
	// method to indicate ehwn a rune was attached
	Attach(*Character, Equippable)

	// Calculates to see if the current rune is active
	IsAttached() bool

	// Calculates to see if the current rune is depleted and item can be deleted
	IsDepleted() bool
}
