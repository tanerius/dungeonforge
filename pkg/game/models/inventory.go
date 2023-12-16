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

// Attachable interface for items that can attach to other items
type Attachable interface {
	Item
	// method to indicate ehwn a rune was attached
	Attach()

	// Calculates to see if the current rune is active
	IsAttached() bool

	// Calculates to see if the current rune is depleted and item can be deleted
	IsDepleted() bool
}
