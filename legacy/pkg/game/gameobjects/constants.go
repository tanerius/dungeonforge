package gameobjects

// Constants defining slots and item types
const (
	// Item worn for head gear
	SlotHead int = 0

	// Item worn as body armour
	SlotArmour int = 1

	// Primary weapon in main hand
	SlotPrimaryW int = 2

	// Secondary weapon or shield hand
	SlotSecondaryW int = 3

	// Boots
	SlotBoots int = 4

	// Ring
	SlotRing int = 5

	// Amulet
	SlotAmulet int = 6

	// Potion buff slot
	SlotBuff1 int = 7

	// Rune buff slot
	SlotBuff2 int = 8

	// Gloves slot
	SlotGloves int = 9

	// Indicator for max
	SlotMaxVal int = 10
)

// constants for rarity
const (
	RarityCommon   = 0 // Gray
	RarityRare     = 1 // Green
	RarityEpic     = 3 // Blue
	RarityImmortal = 4 // YEllow
	RarityMaxVal   = 5
)
