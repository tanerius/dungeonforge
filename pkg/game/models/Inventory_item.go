package models

import "time"

type InventoryItem struct {
	ItemId     string
	BuffIdList string
	IsEquipped bool
	IsBuffed   bool
	BuffedOn   time.Time
}

// Implement Equip required in the Equippable interface
func (i *InventoryItem) Equip(_c *Character) {
	i.IsEquipped = true
}

// Implement Unequip required in the Equippable interface
func (i *InventoryItem) Unequip(_c *Character) {
	i.IsEquipped = false
}
